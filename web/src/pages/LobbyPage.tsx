import { FormEvent, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

import { createRoom, leaveRoom, listRooms, takeSeat } from "../lib/api";
import type { RoomSnapshot } from "../lib/types";
import { viewerSeatKey, viewerTokenKey } from "../lib/viewerSeat";

type LobbyPageProps = {
  navigateToRoom?: (roomId: string) => void;
};

export function LobbyPage({ navigateToRoom }: LobbyPageProps = {}) {
  const [rooms, setRooms] = useState<RoomSnapshot[]>([]);
  const [name, setName] = useState("Table 1");
  const [playerCount, setPlayerCount] = useState(6);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [openingRoomId, setOpeningRoomId] = useState("");
  const [error, setError] = useState("");

  useEffect(() => {
    let cancelled = false;
    async function loadRooms() {
      try {
        const nextRooms = await listRooms();
        if (!cancelled) {
          setRooms(nextRooms);
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "failed to load rooms");
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    void loadRooms();
    return () => {
      cancelled = true;
    };
  }, []);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSubmitting(true);
    setError("");
    try {
      const humanSeat = playerCount - 1;
      const room = await createRoom({
        name,
        smallBlind: 1,
        startingBankroll: 100,
        humanSeat,
        playerCount,
      });
      const session = await takeSeat(room.roomId, room.humanSeat ?? humanSeat);
      if (session.viewerSeat !== undefined) {
        window.localStorage.setItem(
          viewerSeatKey(room.roomId),
          String(session.viewerSeat),
        );
      }
      if (session.viewerToken) {
        window.localStorage.setItem(
          viewerTokenKey(room.roomId),
          session.viewerToken,
        );
      }
      setRooms((current) => [...current, room]);
      setName("Table 1");
      setPlayerCount(6);
      navigateToRoom?.(room.roomId);
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to create room");
    } finally {
      setSubmitting(false);
    }
  }

  async function openRoomAsSpectator(roomId: string) {
    setOpeningRoomId(roomId);
    setError("");
    try {
      const storedViewerToken = window.localStorage.getItem(viewerTokenKey(roomId));
      if (storedViewerToken) {
        await leaveRoom(roomId, storedViewerToken);
      }
      window.localStorage.removeItem(viewerSeatKey(roomId));
      window.localStorage.removeItem(viewerTokenKey(roomId));
      navigateToRoom?.(roomId);
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to open room");
    } finally {
      setOpeningRoomId("");
    }
  }

  return (
    <main className="app-shell">
      <section className="hero-band">
        <div className="hero-copy">
          <span className="eyebrow">Poker Service</span>
          <h1>Poker Control Room</h1>
          <p>
            Create a table, claim the human seat, or monitor the live room feed.
          </p>
        </div>
        <div className="hero-meta">
          <span>
            {loading
              ? "Syncing rooms"
              : `${rooms.length} room${rooms.length === 1 ? "" : "s"}`}
          </span>
          <span>Choose 2-10 players, one human plus AI fillers</span>
        </div>
      </section>

      <section className="work-grid">
        <form className="create-panel" onSubmit={onSubmit}>
          <div className="panel-head">
            <h2>Create room</h2>
            <p>
              Choose the table size. The last seat is reserved for you and the
              other seats are filled by AI.
            </p>
          </div>

          <label className="field">
            <span>Room name</span>
            <input
              value={name}
              onChange={(event) => setName(event.target.value)}
            />
          </label>

          <label className="field">
            <span>Total players</span>
            <select
              value={playerCount}
              onChange={(event) => setPlayerCount(Number(event.target.value))}
            >
              {[2, 3, 4, 5, 6, 7, 8, 9, 10].map((count) => (
                <option key={count} value={count}>
                  {count} players
                </option>
              ))}
            </select>
          </label>

          <button type="submit" disabled={submitting}>
            {submitting ? "Creating..." : "Create Room"}
          </button>

          {error ? <p className="error-text">{error}</p> : null}
        </form>

        <section className="room-strip" aria-label="room list">
          <div className="panel-head">
            <h2>Live rooms</h2>
            <p>
              Each line is a running service room, not a local process snapshot.
            </p>
          </div>

          <div className="room-list">
            {rooms.map((room) => {
              return (
                <article key={room.roomId} className="room-row">
                  <div className="room-row-main">
                    <div className="room-row-title">
                      <h3>{room.roomName}</h3>
                      <span className="status-pill">{room.status}</span>
                    </div>
                    <button
                      className="room-open-link"
                      disabled={openingRoomId === room.roomId}
                      onClick={() => void openRoomAsSpectator(room.roomId)}
                      type="button"
                    >
                      {openingRoomId === room.roomId
                        ? "Opening..."
                        : "Spectate room"}
                    </button>
                  </div>
                  <div className="room-row-meta">
                    <span>Room {room.roomId}</span>
                    <span>Players {room.playerCount ?? 6}</span>
                    <span>Blind {room.smallBlind}</span>
                    <span>Hand {room.handNumber}</span>
                    <span>Human seat {room.humanSeat ?? 5}</span>
                  </div>
                </article>
              );
            })}

            {!loading && rooms.length === 0 ? (
              <p className="empty-state">
                No rooms yet. Create the first table.
              </p>
            ) : null}
          </div>
        </section>
      </section>
    </main>
  );
}

export function LobbyRoute() {
  const navigate = useNavigate();
  return (
    <LobbyPage navigateToRoom={(roomId) => navigate(`/rooms/${roomId}`)} />
  );
}
