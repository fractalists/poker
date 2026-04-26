import { useEffect, useState } from "react";

import type { ActionSubmission, PendingAction, SeatSnapshot } from "../lib/types";
import { CardFace } from "./CardFace";

type ActionBarProps = {
  roomId: string;
  pot?: number;
  pendingAction: PendingAction;
  playerSeat?: SeatSnapshot;
  busy?: boolean;
  onSubmit: (action: ActionSubmission) => void | Promise<void>;
};

function getMinimumBet(pendingAction: PendingAction) {
  return pendingAction.minBetAmount ?? Math.max(pendingAction.minAmount + 1, 1);
}

function clampBetAmount(amount: number, pendingAction: PendingAction) {
  const minimumBet = getMinimumBet(pendingAction);
  if (!Number.isFinite(amount)) {
    return minimumBet;
  }
  return Math.min(pendingAction.maxAmount, Math.max(minimumBet, Math.ceil(amount)));
}

function formatBetRange(pendingAction: PendingAction) {
  const minimumBet = getMinimumBet(pendingAction);
  const actionLabel = pendingAction.minAmount > 0 ? "Raise" : "Bet";
  if (minimumBet >= pendingAction.maxAmount) {
    return `${actionLabel} ${pendingAction.maxAmount}`;
  }
  return `${actionLabel} ${minimumBet} to ${pendingAction.maxAmount}`;
}

function getShortcutBetAmount(pot: number, ratio: number, pendingAction: PendingAction) {
  const callAmount = pendingAction.minAmount;
  const potAfterCall = pot + callAmount;
  return clampBetAmount(callAmount + ratio * potAfterCall, pendingAction);
}

function formatPotOdds(pot: number, callAmount: number) {
  if (callAmount <= 0) {
    return "Free";
  }

  const potAfterCall = pot + callAmount;
  if (potAfterCall <= 0) {
    return "--";
  }

  return `${Math.round((callAmount / potAfterCall) * 100)}%`;
}

type ShortcutKey = "quarter" | "half" | "pot";

const betShortcuts: Array<{ key: ShortcutKey; label: string; ratio: number }> = [
  { key: "quarter", label: "1/4 Pot", ratio: 0.25 },
  { key: "half", label: "1/2 Pot", ratio: 0.5 },
  { key: "pot", label: "1 Pot", ratio: 1 },
];

export function ActionBar({
  roomId,
  pot = 0,
  pendingAction,
  playerSeat,
  busy = false,
  onSubmit,
}: ActionBarProps) {
  const [betAmount, setBetAmount] = useState("");
  const [isBetting, setIsBetting] = useState(false);
  const [selectedShortcut, setSelectedShortcut] = useState<ShortcutKey | null>(null);
  const [remainingSeconds, setRemainingSeconds] = useState<number | null>(() =>
    pendingAction.expiresAt
      ? Math.max(0, Math.ceil((pendingAction.expiresAt - Date.now()) / 1000))
      : null,
  );
  const isCallEquivalentToAllIn = pendingAction.canAllIn && pendingAction.minAmount >= pendingAction.maxAmount;
  const canShowCall =
    pendingAction.canCheck || (pendingAction.canCall && pendingAction.minAmount <= pendingAction.maxAmount && !isCallEquivalentToAllIn);
  const callLabel = pendingAction.canCheck ? "Check" : `Call ${pendingAction.minAmount}`;
  const betLabel = formatBetRange(pendingAction);
  const isRaise = pendingAction.minAmount > 0;
  const amountLabel = isRaise ? "Raise to amount" : "Bet amount";
  const confirmLabel = isRaise ? "Confirm raise" : "Confirm bet";
  const minimumBet = getMinimumBet(pendingAction);
  const isBetEquivalentToAllIn = pendingAction.canAllIn && minimumBet >= pendingAction.maxAmount;
  const canShowBet = pendingAction.canBet && !isBetEquivalentToAllIn;
  const stackAmount = pendingAction.maxAmount;
  const callAmount = pendingAction.canCheck
    ? 0
    : Math.min(Math.max(pendingAction.minAmount, 0), stackAmount);
  const afterCallStack = Math.max(0, stackAmount - callAmount);
  const potOdds = formatPotOdds(pot, callAmount);
  const decisionCards = playerSeat?.cards?.length ? playerSeat.cards : ["--", "--"];

  useEffect(() => {
    setBetAmount("");
    setIsBetting(false);
    setSelectedShortcut(null);
    setRemainingSeconds(
      pendingAction.expiresAt
        ? Math.max(0, Math.ceil((pendingAction.expiresAt - Date.now()) / 1000))
        : null,
    );
  }, [roomId, pendingAction.token]);

  useEffect(() => {
    if (!pendingAction.expiresAt) {
      setRemainingSeconds(null);
      return;
    }

    function updateRemainingSeconds() {
      setRemainingSeconds(
        Math.max(0, Math.ceil(((pendingAction.expiresAt ?? 0) - Date.now()) / 1000)),
      );
    }

    updateRemainingSeconds();
    const timer = window.setInterval(updateRemainingSeconds, 250);
    return () => window.clearInterval(timer);
  }, [pendingAction.expiresAt, pendingAction.token]);

  function applyBetShortcut(shortcutKey: ShortcutKey, ratio: number) {
    setBetAmount(String(getShortcutBetAmount(pot, ratio, pendingAction)));
    setSelectedShortcut(shortcutKey);
  }

  function submitBet() {
    return onSubmit({
      token: pendingAction.token,
      actionType: "BET",
      amount: clampBetAmount(Number(betAmount), pendingAction),
    });
  }

  return (
    <section className="action-bar" aria-label="action controls">
      <div className="action-bar-top">
        <div>
          <span className="eyebrow">Your Turn</span>
        </div>
        {remainingSeconds !== null ? (
          <strong className="turn-countdown">{remainingSeconds}s</strong>
        ) : null}
      </div>

      <div
        className={[
          "decision-summary",
          playerSeat ? "" : "decision-summary--metrics-only",
        ]
          .filter(Boolean)
          .join(" ")}
      >
        {playerSeat ? (
          <div className="decision-hand">
            <span className="decision-label">Your hand</span>
            <div className="decision-cards" aria-label="your hand cards">
              {decisionCards.map((card, index) => (
                <span className="decision-card" key={`${card}-${index}`}>
                  <CardFace card={card} />
                </span>
              ))}
            </div>
          </div>
        ) : null}

        <div className="decision-metrics">
          <span className="decision-metric decision-metric--pot">
            <span className="decision-label">Pot</span>
            <strong>{pot}</strong>
          </span>
          <span className="decision-metric">
            <span className="decision-label">Stack</span>
            <strong>{stackAmount}</strong>
          </span>
          <span className="decision-metric">
            <span className="decision-label">After call</span>
            <strong>{afterCallStack}</strong>
          </span>
          <span className="decision-metric decision-metric--odds">
            <span className="decision-label">Pot odds</span>
            <strong>{potOdds}</strong>
          </span>
        </div>
      </div>

      <div className="action-buttons">
        {pendingAction.canFold ? (
          <button className="action-button action-button--fold" disabled={busy} onClick={() => onSubmit({ token: pendingAction.token, actionType: "FOLD", amount: 0 })} type="button">
            Fold
          </button>
        ) : null}
        {canShowCall ? (
          <button
            className="action-button action-button--call primary"
            disabled={busy}
            onClick={() =>
              onSubmit({
                token: pendingAction.token,
                actionType: "CALL",
                amount: pendingAction.canCheck ? 0 : pendingAction.minAmount,
              })
            }
            type="button"
          >
            {callLabel}
          </button>
        ) : null}
        {canShowBet ? (
          <button
            className={[
              "action-button",
              "action-button--raise",
              isBetting ? "is-open" : "",
            ]
              .filter(Boolean)
              .join(" ")}
            disabled={busy}
            onClick={() => setIsBetting((current) => !current)}
            type="button"
          >
            {betLabel}
          </button>
        ) : null}
        {pendingAction.canAllIn ? (
          <button
            className="action-button action-button--all-in"
            disabled={busy}
            onClick={() =>
              onSubmit({
                token: pendingAction.token,
                actionType: "ALL_IN",
                amount: pendingAction.maxAmount,
              })
            }
            type="button"
          >
            {`All-in ${pendingAction.maxAmount}`}
          </button>
        ) : null}
      </div>

      {canShowBet && isBetting ? (
        <div className="bet-panel">
          <label className="action-amount">
            <span>{amountLabel}</span>
            <input
              aria-label={amountLabel}
              type="number"
              min={minimumBet}
              max={pendingAction.maxAmount}
              value={betAmount}
              onChange={(event) => {
                setBetAmount(event.target.value);
                setSelectedShortcut(null);
              }}
            />
          </label>

          <div className="bet-shortcuts">
            {betShortcuts.map((shortcut) => (
              <button
                aria-pressed={selectedShortcut === shortcut.key}
                className={selectedShortcut === shortcut.key ? "is-selected" : undefined}
                disabled={busy}
                key={shortcut.key}
                onClick={() => applyBetShortcut(shortcut.key, shortcut.ratio)}
                type="button"
              >
                {shortcut.label}
              </button>
            ))}
            <button className="primary" disabled={busy} onClick={() => void submitBet()} type="button">
              {confirmLabel}
            </button>
          </div>
        </div>
      ) : null}
    </section>
  );
}
