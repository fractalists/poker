import { useEffect, useRef, useState, type CSSProperties } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";

import { ActionBar } from "../components/ActionBar";
import { CardFace } from "../components/CardFace";
import { TableSeat, type SeatActionView } from "../components/TableSeat";
import {
  getRoom,
  leaveRoom,
  startHand,
  submitAction,
  takeSeat,
} from "../lib/api";
import { subscribeRoom } from "../lib/socket";
import type {
  ActionSubmission,
  RoomSnapshot,
  SeatSnapshot,
} from "../lib/types";
import { viewerSeatKey, viewerTokenKey } from "../lib/viewerSeat";
import { buildOrbitLayout } from "./roomLayout";

type RoomPageProps = {
  snapshot: RoomSnapshot;
  busy?: boolean;
  error?: string;
  onAction: (action: ActionSubmission) => Promise<void>;
  onStartHand: () => Promise<void>;
  onTakeSeat: () => Promise<void>;
  onLeave?: () => Promise<void>;
  onSpectate?: () => Promise<void>;
};

function formatNetChange(delta: number) {
  if (delta > 0) {
    return `+${delta}`;
  }
  return `${delta}`;
}

function buildSettlementHeadline(seats: SeatSnapshot[]) {
  const winners = seats.filter((seat) => seat.isWinner);
  if (winners.length === 0) {
    return "Hand complete";
  }
  if (winners.length === 1) {
    const winner = winners[0];
    return `${winner.name} wins ${formatNetChange(winner.netChange ?? 0)}`;
  }
  return `${winners.map((seat) => seat.name).join(" + ")} split the pot`;
}

function buildSettlementDetail(seats: SeatSnapshot[]) {
  const winners = seats.filter((seat) => seat.isWinner);
  if (winners.length === 1 && winners[0].bestHand) {
    return `Winning hand: ${winners[0].bestHand}`;
  }
  if (winners.length > 1) {
    const sharedHand = winners[0]?.bestHand;
    if (sharedHand && winners.every((seat) => seat.bestHand === sharedHand)) {
      return `Winning hand: ${sharedHand}`;
    }
  }
  return "Review the final hands and chip swings below before starting the next hand.";
}

function formatRoomStatus(status: string) {
  return status.replace(/_/g, " ");
}

function formatSeatActionLabel(actionType?: string, amount?: number) {
  switch (actionType) {
    case "CALL":
      return amount && amount > 0 ? `Call ${amount}` : "Check";
    case "BET":
      return amount !== undefined ? `Bet ${amount}` : "Bet";
    case "FOLD":
      return "Folded";
    case "ALL_IN":
      return amount !== undefined ? `All-in ${amount}` : "All-in";
    default:
      return "";
  }
}

function mapActionTone(actionType?: string): SeatActionView["tone"] {
  switch (actionType) {
    case "FOLD":
      return "fold";
    case "ALL_IN":
      return "all-in";
    case "BET":
      return "bet";
    default:
      return "call";
  }
}

function buildSeatActionMap(snapshot: RoomSnapshot) {
  const currentRound = snapshot.round;
  const currentHand = snapshot.handNumber;
  const events = snapshot.events ?? [];

  if (!currentRound) {
    return new Map<number, SeatActionView>();
  }

  let startIndex = 0;
  for (let index = 0; index < events.length; index += 1) {
    const event = events[index];
    if (
      event.kind === "round_start" &&
      event.handNumber === currentHand &&
      event.round === currentRound
    ) {
      startIndex = index + 1;
    }
  }

  const actions = new Map<number, SeatActionView>();
  for (let index = startIndex; index < events.length; index += 1) {
    const event = events[index];
    if (
      event.kind !== "player_action" ||
      event.handNumber !== currentHand ||
      event.round !== currentRound ||
      event.seatIndex === undefined
    ) {
      continue;
    }

    const label = formatSeatActionLabel(event.actionType, event.amount);
    if (!label) {
      continue;
    }

    actions.set(event.seatIndex, {
      label,
      tone: mapActionTone(event.actionType),
      stamp: `${currentHand}-${currentRound}-${index}`,
    });
  }

  return actions;
}

export function RoomPage({
  snapshot,
  busy = false,
  error = "",
  onAction,
  onStartHand,
  onTakeSeat,
  onLeave,
  onSpectate,
}: RoomPageProps) {
  const leaveConfirmRef = useRef<HTMLDivElement | null>(null);
  const [showLeaveConfirm, setShowLeaveConfirm] = useState(false);
  const boardCards = snapshot.boardCards ?? ["**", "**", "**", "**", "**"];
  const dealtBoardCards = snapshot.boardCards ?? [];
  const seats = snapshot.seats ?? [];
  const events = snapshot.events ?? [];
  const seatActions = buildSeatActionMap(snapshot);
  const orbitPlayerCount = Math.max(
    snapshot.playerCount ?? 0,
    seats.length,
    snapshot.humanSeat !== undefined ? snapshot.humanSeat + 1 : 0,
  );
  const orbitLayout = buildOrbitLayout(
    Math.max(2, orbitPlayerCount || 2),
    snapshot.humanSeat,
  );
  const { spec: orbitSpec, positions: orbitPositions } = orbitLayout;
  const orbitStyle = {
    "--orbit-min-height": orbitSpec.minHeight,
    "--orbit-seat-width": orbitSpec.seatWidth,
    "--orbit-seat-min-height": orbitSpec.seatMinHeight,
    "--orbit-card-scale": String(orbitSpec.cardScale),
    "--orbit-board-top-padding": orbitSpec.boardTopPadding,
    "--orbit-board-side-padding": orbitSpec.boardSidePadding,
    "--orbit-board-bottom-padding": orbitSpec.boardBottomPadding,
    "--orbit-board-width": orbitSpec.boardWidth,
    "--orbit-stat-width": orbitSpec.statWidth,
    "--orbit-board-gap": orbitSpec.boardGap,
  } as CSSProperties;
  const isPlayerView = snapshot.viewerRole === "player";
  const isHandFinished = snapshot.status === "hand_finished";
  const canStartHand =
    snapshot.status === "waiting" || snapshot.status === "hand_finished";
  const startHandLabel = canStartHand ? "Start hand" : "Hand in progress";
  const isOpeningState =
    snapshot.status === "waiting" &&
    !snapshot.pendingAction &&
    !isHandFinished &&
    seats.length === 0 &&
    dealtBoardCards.length === 0;
  const settlementSeats = isHandFinished
    ? [...seats].sort(
        (left, right) =>
          (right.netChange ?? 0) - (left.netChange ?? 0) ||
          left.index - right.index,
      )
    : [];
  const tableStats = [
    { label: "Round", value: snapshot.round ?? "WAITING" },
    { label: "Pot", value: String(snapshot.pot ?? 0) },
    { label: "Current", value: String(snapshot.currentAmount ?? 0) },
  ];
  const metaChips = [
    {
      label: formatRoomStatus(snapshot.status),
      className: "room-chip room-chip--status",
    },
    ...(isPlayerView
      ? []
      : [{ label: "spectator view", className: "room-chip room-chip--view" }]),
    {
      label: `hand ${snapshot.handNumber}`,
      className: "room-chip room-chip--hand",
    },
  ];

  useEffect(() => {
    if (!showLeaveConfirm) {
      return;
    }

    function handlePointerDown(event: MouseEvent) {
      if (!leaveConfirmRef.current?.contains(event.target as Node)) {
        setShowLeaveConfirm(false);
      }
    }

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === "Escape") {
        setShowLeaveConfirm(false);
      }
    }

    document.addEventListener("mousedown", handlePointerDown);
    document.addEventListener("keydown", handleKeyDown);

    return () => {
      document.removeEventListener("mousedown", handlePointerDown);
      document.removeEventListener("keydown", handleKeyDown);
    };
  }, [showLeaveConfirm]);

  return (
    <main className="room-shell room-shell--fixed">
      <header className="room-topbar">
        <div className="room-title-block">
          <div className="room-title-row">
            <span className="eyebrow">Live Room</span>
            <h1>{snapshot.roomName}</h1>
          </div>
        </div>
        <div className="room-topbar-actions">
          <div className="room-meta-strip">
            {metaChips.map((chip) => (
              <span className={chip.className} key={chip.label}>
                {chip.label}
              </span>
            ))}
          </div>
          <div className="room-exit-wrap" ref={leaveConfirmRef}>
            <button
              aria-controls="leave-room-confirmation"
              aria-expanded={showLeaveConfirm}
              aria-haspopup="dialog"
              className={[
                "room-leave-trigger",
                showLeaveConfirm ? "is-open" : "",
              ]
                .filter(Boolean)
                .join(" ")}
              onClick={() => setShowLeaveConfirm((current) => !current)}
              type="button"
            >
              Back to lobby
            </button>
            {showLeaveConfirm ? (
              <div
                aria-label="Leave room confirmation"
                className="room-exit-confirm"
                id="leave-room-confirmation"
                role="dialog"
              >
                <p className="room-exit-confirm-title">Leave this table?</p>
                <p className="room-exit-confirm-copy">
                  You will return to the lobby and stop following the current
                  hand from this page.
                </p>
                <div className="room-exit-confirm-actions">
                  <button
                    className="room-exit-stay"
                    onClick={() => setShowLeaveConfirm(false)}
                    type="button"
                  >
                    Stay here
                  </button>
                  {onLeave ? (
                    <Link
                      className="room-exit-link"
                      onClick={(event) => {
                        event.preventDefault();
                        void onLeave();
                      }}
                      to="/"
                    >
                      Leave room
                    </Link>
                  ) : (
                    <a className="room-exit-link" href="/">
                      Leave room
                    </a>
                  )}
                </div>
              </div>
            ) : null}
          </div>
        </div>
      </header>

      <section className="room-grid">
        <section className="table-stage">
          <div className="table-felt">
            {isOpeningState ? (
              <div className="table-empty-state">
                <div aria-hidden="true" className="table-empty-orbit">
                  <span className="table-empty-puck" />
                </div>
                <div className="table-empty-copy">
                  <span className="eyebrow">Table Ready</span>
                  <h2>Waiting for the opening deal</h2>
                  <p>
                    Start a hand to bring players onto the felt. Community cards
                    and live betting state will appear once the action begins.
                  </p>
                </div>
              </div>
            ) : (
              <>
                <div className="table-live-layout" style={orbitStyle}>
                  <div className="board-strip">
                    <div className="board-cluster">
                      <div className="board-meta" aria-label="table stats">
                        {tableStats.map((stat) => (
                          <div className="table-stat-row" key={stat.label}>
                            <div className="table-stat-copy">
                              <span className="table-stat-label">
                                {stat.label}
                              </span>
                              <strong className="table-stat-value">
                                {stat.value}
                              </strong>
                            </div>
                          </div>
                        ))}
                      </div>
                      <div className="board-centerpiece">
                        <div className="board-badge">
                          <p className="table-note table-note--board">
                            Community cards
                          </p>
                        </div>
                        <div className="board-cards-shell">
                          <div className="board-cards">
                            {boardCards.map((card, index) => (
                              <div
                                className="board-card"
                                key={`${card}-${index}`}
                              >
                                <CardFace card={card} />
                              </div>
                            ))}
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>

                  <div
                    className={`seat-orbit seat-orbit--${Math.max(
                      2,
                      orbitPlayerCount || 2,
                    )}`}
                  >
                    {seats.length === 0 ? (
                      <div className="table-placeholder">
                        Start a hand to deal players into the table.
                      </div>
                    ) : null}
                    {seats.map((seat) => {
                      const slot = orbitPositions.get(seat.index);
                      return (
                        <div
                          className={[
                            "seat-slot",
                            slot ? `seat-slot--${slot.slot}` : "",
                          ]
                            .filter(Boolean)
                            .join(" ")}
                          key={seat.index}
                          style={
                            slot
                              ? ({
                                  left: `${slot.x}%`,
                                  top: `${slot.y}%`,
                                } as CSSProperties)
                              : undefined
                          }
                        >
                          <TableSeat
                            recentAction={seatActions.get(seat.index)}
                            seat={seat}
                            viewerSeat={
                              isPlayerView ? snapshot.humanSeat : undefined
                            }
                          />
                        </div>
                      );
                    })}
                  </div>
                </div>
              </>
            )}
          </div>
        </section>

        <aside className="side-panel side-panel--scroll">
          <section className="control-stack">
            <button
              className="primary"
              disabled={busy || !canStartHand}
              onClick={() => void onStartHand()}
              type="button"
            >
              {startHandLabel}
            </button>
            {snapshot.viewerRole === "player" ? null : (
              <button
                disabled={busy}
                onClick={() => void onTakeSeat()}
                type="button"
              >
                Take human seat
              </button>
            )}
          </section>

          {isPlayerView && snapshot.pendingAction ? (
            <ActionBar
              roomId={snapshot.roomId}
              pot={snapshot.pot ?? 0}
              pendingAction={snapshot.pendingAction}
              busy={busy}
              onSubmit={onAction}
            />
          ) : null}

          {!isPlayerView && !isHandFinished ? (
            <section className="viewer-note">
              <span className="eyebrow">View Mode</span>
              <p>
                Your cards are hidden because this room is currently in
                spectator mode.
              </p>
            </section>
          ) : null}

          {error ? <p className="error-text">{error}</p> : null}

          <div className="room-history-stack room-history-stack--scroll">
            {isHandFinished ? (
              <section className="settlement-panel">
                <span className="eyebrow">Hand Review</span>
                <h2>{buildSettlementHeadline(settlementSeats)}</h2>
                <p>{buildSettlementDetail(settlementSeats)}</p>
                <ul className="settlement-list">
                  {settlementSeats.map((seat) => (
                    <li key={`settlement-${seat.index}`}>
                      <div>
                        <strong>{seat.name}</strong>
                        {seat.bestHand ? (
                          <span>{seat.bestHand}</span>
                        ) : (
                          <span>
                            {seat.isWinner
                              ? "Won without showdown"
                              : "No revealed hand"}
                          </span>
                        )}
                      </div>
                      <span
                        className={[
                          "settlement-delta",
                          (seat.netChange ?? 0) > 0 ? "is-up" : "",
                          (seat.netChange ?? 0) < 0 ? "is-down" : "",
                        ]
                          .filter(Boolean)
                          .join(" ")}
                      >
                        {formatNetChange(seat.netChange ?? 0)}
                      </span>
                    </li>
                  ))}
                </ul>
              </section>
            ) : null}

            <section className="room-feed-panel">
              <div className="panel-head">
                <h2>Table feed</h2>
                <p>Actions and lifecycle events from the room runtime.</p>
              </div>
              <ul className="event-feed">
                {events.length === 0 ? <li>No events yet.</li> : null}
                {events
                  .slice(-8)
                  .reverse()
                  .map((event, index) => (
                    <li key={`${event.kind}-${index}`}>{event.message}</li>
                  ))}
              </ul>
            </section>
          </div>
        </aside>
      </section>
    </main>
  );
}

export function RoomRoute() {
  const navigate = useNavigate();
  const params = useParams();
  const roomId = params.roomId ?? "";
  const [viewerSeat, setViewerSeat] = useState<number | undefined>(() => {
    const stored = roomId
      ? window.localStorage.getItem(viewerSeatKey(roomId))
      : null;
    return stored === null ? undefined : Number(stored);
  });
  const [viewerToken, setViewerToken] = useState<string | undefined>(() => {
    const stored = roomId
      ? window.localStorage.getItem(viewerTokenKey(roomId))
      : null;
    return stored ?? undefined;
  });
  const [snapshot, setSnapshot] = useState<RoomSnapshot | null>(null);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!roomId) {
      return;
    }

    let cancelled = false;
    void getRoom(roomId, viewerSeat, viewerToken)
      .then((nextSnapshot) => {
        if (!cancelled) {
          setSnapshot(nextSnapshot);
        }
      })
      .catch((err) => {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "failed to load room");
        }
      });

    return () => {
      cancelled = true;
    };
  }, [roomId, viewerSeat, viewerToken]);

  useEffect(() => {
    if (!roomId) {
      return;
    }

    return subscribeRoom(
      roomId,
      viewerSeat,
      viewerToken,
      (nextSnapshot) => {
        setSnapshot(nextSnapshot);
        setError("");
      },
      (message) => setError(message),
    );
  }, [roomId, viewerSeat, viewerToken]);

  async function withBusy(work: () => Promise<void>) {
    setBusy(true);
    setError("");
    try {
      await work();
    } catch (err) {
      setError(err instanceof Error ? err.message : "request failed");
    } finally {
      setBusy(false);
    }
  }

  if (!roomId) {
    return (
      <main className="app-shell">
        <p className="error-text">Missing room id.</p>
      </main>
    );
  }

  if (snapshot === null && error) {
    return (
      <main className="app-shell">
        <section className="create-panel">
          <span className="eyebrow">Room unavailable</span>
          <div className="panel-head">
            <h2>We couldn&apos;t open this table.</h2>
            <p>{error}</p>
          </div>
          <Link className="room-open-link" to="/">
            Back to lobby
          </Link>
        </section>
      </main>
    );
  }

  if (snapshot === null) {
    return (
      <main className="app-shell">
        <p className="table-note">Loading room...</p>
      </main>
    );
  }

  return (
    <RoomPage
      snapshot={snapshot}
      busy={busy}
      error={error}
      onAction={(action) =>
        withBusy(() => submitAction(roomId, action, viewerToken))
      }
      onStartHand={() => withBusy(() => startHand(roomId))}
      onTakeSeat={() =>
        withBusy(async () => {
          const session = await takeSeat(roomId, snapshot.humanSeat ?? 5);
          if (session.viewerSeat !== undefined) {
            window.localStorage.setItem(
              viewerSeatKey(roomId),
              String(session.viewerSeat),
            );
            setViewerSeat(session.viewerSeat);
          }
          if (session.viewerToken) {
            window.localStorage.setItem(
              viewerTokenKey(roomId),
              session.viewerToken,
            );
            setViewerToken(session.viewerToken);
          }
        })
      }
      onLeave={() =>
        withBusy(async () => {
          if (viewerToken) {
            await leaveRoom(roomId, viewerToken);
          }
          window.localStorage.removeItem(viewerSeatKey(roomId));
          window.localStorage.removeItem(viewerTokenKey(roomId));
          setViewerSeat(undefined);
          setViewerToken(undefined);
          navigate("/");
        })
      }
      onSpectate={() =>
        withBusy(async () => {
          await leaveRoom(roomId, viewerToken);
          window.localStorage.removeItem(viewerSeatKey(roomId));
          window.localStorage.removeItem(viewerTokenKey(roomId));
          setViewerSeat(undefined);
          setViewerToken(undefined);
        })
      }
    />
  );
}
