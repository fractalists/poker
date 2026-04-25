import type { SeatSnapshot } from "../lib/types";
import { CardFace } from "./CardFace";

export type SeatActionView = {
  label: string;
  tone: "call" | "bet" | "fold" | "all-in" | "out";
  stamp: string;
};

type TableSeatProps = {
  seat: SeatSnapshot;
  viewerSeat?: number;
  recentAction?: SeatActionView;
  showSettlementEffects?: boolean;
};

function formatNetChange(delta: number) {
  if (delta > 0) {
    return `+${delta}`;
  }
  return `${delta}`;
}

export function TableSeat({
  seat,
  viewerSeat,
  recentAction,
  showSettlementEffects = false,
}: TableSeatProps) {
  const cards = seat.cards ?? [];
  const hasOutcome = Boolean(seat.bestHand);
  const isViewerSeat = viewerSeat === seat.index;
  const isAllInSeat = seat.status === "ALLIN" || seat.status === "ALL_IN";
  const isEliminatedSeat = seat.status === "OUT" && seat.bankroll === 0;
  const isFoldedSeat = seat.status === "OUT" && !isEliminatedSeat;
  const isSettlementWinner = showSettlementEffects && Boolean(seat.isWinner);
  const isSettlementLoser =
    showSettlementEffects && !seat.isWinner && (seat.netChange ?? 0) < 0;
  const headerBadge = recentAction
    ? recentAction
    : seat.isTurn
      ? { label: "To act", tone: "call" as const, stamp: `seat-${seat.index}-turn` }
      : isAllInSeat
        ? { label: "All-in", tone: "all-in" as const, stamp: `seat-${seat.index}-all-in` }
    : isEliminatedSeat
      ? { label: "Busted", tone: "out" as const, stamp: `seat-${seat.index}-out` }
      : null;
  const className = [
    "table-seat",
    headerBadge ? "has-recent-action" : "",
    recentAction ? `has-recent-action--${recentAction.tone}` : "",
    seat.isTurn ? "is-turn" : "",
    isAllInSeat ? "is-all-in" : "",
    isFoldedSeat ? "is-folded" : "",
    isEliminatedSeat ? "is-eliminated" : "",
    isSettlementWinner ? "is-settlement-winner" : "",
    isSettlementLoser ? "is-settlement-loser" : "",
    seat.isWinner ? "is-winner" : "",
    isViewerSeat ? "is-viewer" : "",
  ]
    .filter(Boolean)
    .join(" ");

  return (
    <article className={className}>
      <header className="seat-header">
        <div className="seat-identity">
          <div className="seat-name-row">
            <h3>{seat.name}</h3>
            {seat.position ? (
              <span className="seat-position-badge">{seat.position}</span>
            ) : null}
            {isViewerSeat ? <span className="seat-player-badge">You</span> : null}
          </div>
        </div>
        <div className="seat-header-meta">
          {headerBadge ? (
            <span
              className={`seat-action-pill tone-${headerBadge.tone}`}
              key={headerBadge.stamp}
            >
              {headerBadge.label}
            </span>
          ) : null}
        </div>
      </header>

      <div className="seat-body">
        <div className="seat-cards">
          {(cards.length > 0 ? cards : ["--", "--"]).map((card, index) => (
            <div className="seat-card" key={`${seat.index}-${index}-${card}`}>
              <CardFace card={card} />
            </div>
          ))}
        </div>

        <div className="seat-meta-row">
          <div className="seat-outcome">
            {hasOutcome ? (
              <span className="seat-result-label">{seat.bestHand}</span>
            ) : null}
          </div>
        </div>
      </div>

      <footer className="seat-footer">
        <div className="seat-bankroll-group">
          <span className="seat-stack-pill seat-stack-pill--bankroll">
            <span className="seat-stack-label">Bankroll</span>
            <span className="seat-stack-value">{seat.bankroll}</span>
            {seat.netChange !== undefined ? (
              <span
                className={[
                  "seat-result-pill",
                  seat.netChange > 0 ? "is-up" : "",
                  seat.netChange < 0 ? "is-down" : "",
                ]
                  .filter(Boolean)
                  .join(" ")}
              >
                {formatNetChange(seat.netChange)}
              </span>
            ) : null}
          </span>
        </div>
        <span className="seat-stack-pill seat-stack-pill--pot">
          <span className="seat-stack-label">In pot</span>
          <span className="seat-stack-value">{seat.inPotAmount}</span>
        </span>
      </footer>
    </article>
  );
}
