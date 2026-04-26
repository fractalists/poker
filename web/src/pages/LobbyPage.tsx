import { FormEvent, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

import { createRoom, leaveRoom, listRooms, takeSeat } from "../lib/api";
import { subscribeRooms, type RoomSocketStatus } from "../lib/socket";
import type { RoomSnapshot } from "../lib/types";
import { viewerSeatKey, viewerTokenKey } from "../lib/viewerSeat";

type LobbyPageProps = {
  navigateToRoom?: (roomId: string) => void;
};

function formatRoomAIStyle(style?: string) {
  const normalized = style?.trim().toLowerCase();
  if (!normalized || normalized === "mixed" || normalized === "random") {
    return "mixed";
  }
  return normalized;
}

export function LobbyPage({ navigateToRoom }: LobbyPageProps = {}) {
  const [rooms, setRooms] = useState<RoomSnapshot[]>([]);
  const [name, setName] = useState("Table 1");
  const [playerCount, setPlayerCount] = useState(6);
  const [aiStyle, setAIStyle] = useState("random");
  const [loading, setLoading] = useState(true);
  const [roomSyncStatus, setRoomSyncStatus] = useState<
    RoomSocketStatus | "connecting"
  >("connecting");
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

  useEffect(() => {
    return subscribeRooms(
      (nextRooms) => {
        setRooms(nextRooms);
        setError("");
        setLoading(false);
      },
      (message) => setError(message),
      (status) => setRoomSyncStatus(status),
    );
  }, []);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSubmitting(true);
    setError("");
    try {
      const humanSeat = Math.floor(Math.random() * playerCount);
      const room = await createRoom({
        name,
        smallBlind: 1,
        startingBankroll: 100,
        humanSeat,
        playerCount,
        aiStyle,
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
      setAIStyle("random");
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
          <span>
            {roomSyncStatus === "live" ? "Live lobby" : roomSyncStatus}
          </span>
          <span>Choose 2-10 players, one human plus AI fillers</span>
        </div>
      </section>

      <section className="work-grid">
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
                    <span>AI {formatRoomAIStyle(room.aiStyle)}</span>
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

        <form className="create-panel" onSubmit={onSubmit}>
          <div className="panel-head">
            <h2>Create room</h2>
            <p>
              Choose the table size. Your seat is assigned randomly and the
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

          <label className="field">
            <span>AI style</span>
            <select
              value={aiStyle}
              onChange={(event) => setAIStyle(event.target.value)}
            >
              <option value="random">Mixed AI</option>
              <option value="gto">GTO-inspired</option>
              <option value="smart">Smart odds</option>
              <option value="conservative">Conservative</option>
              <option value="aggressive">Aggressive</option>
            </select>
          </label>

          <button type="submit" disabled={submitting}>
            {submitting ? "Creating..." : "Create Room"}
          </button>

          {error ? <p className="error-text">{error}</p> : null}
        </form>
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
