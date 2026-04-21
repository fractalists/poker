import { useEffect, useState } from "react";

import type { ActionSubmission, PendingAction } from "../lib/types";

type ActionBarProps = {
  roomId: string;
  pot?: number;
  pendingAction: PendingAction;
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
  if (minimumBet >= pendingAction.maxAmount) {
    return `Bet ${pendingAction.maxAmount}`;
  }
  return `Bet ${minimumBet}~${pendingAction.maxAmount}`;
}

function getShortcutBetAmount(pot: number, ratio: number, pendingAction: PendingAction) {
  const callAmount = pendingAction.minAmount;
  const potAfterCall = pot + callAmount;
  return clampBetAmount(callAmount + ratio * potAfterCall, pendingAction);
}

type ShortcutKey = "quarter" | "half" | "pot";

const betShortcuts: Array<{ key: ShortcutKey; label: string; ratio: number }> = [
  { key: "quarter", label: "1/4 Pot", ratio: 0.25 },
  { key: "half", label: "1/2 Pot", ratio: 0.5 },
  { key: "pot", label: "1 Pot", ratio: 1 },
];

export function ActionBar({ roomId, pot = 0, pendingAction, busy = false, onSubmit }: ActionBarProps) {
  const [betAmount, setBetAmount] = useState("");
  const [isBetting, setIsBetting] = useState(false);
  const [selectedShortcut, setSelectedShortcut] = useState<ShortcutKey | null>(null);
  const isCallEquivalentToAllIn = pendingAction.canAllIn && pendingAction.minAmount >= pendingAction.maxAmount;
  const canShowCall =
    pendingAction.canCheck || (pendingAction.canCall && pendingAction.minAmount <= pendingAction.maxAmount && !isCallEquivalentToAllIn);
  const callLabel = pendingAction.canCheck ? "Check" : `Call ${pendingAction.minAmount}`;
  const betLabel = formatBetRange(pendingAction);

  useEffect(() => {
    setBetAmount("");
    setIsBetting(false);
    setSelectedShortcut(null);
  }, [roomId, pendingAction.token]);

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
      </div>

      <div className="action-buttons">
        {pendingAction.canFold ? (
          <button disabled={busy} onClick={() => onSubmit({ token: pendingAction.token, actionType: "FOLD", amount: 0 })} type="button">
            Fold
          </button>
        ) : null}
        {canShowCall ? (
          <button
            className="primary"
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
        {pendingAction.canBet ? (
          <button
            className={isBetting ? "is-open" : undefined}
            disabled={busy}
            onClick={() => setIsBetting((current) => !current)}
            type="button"
          >
            {betLabel}
          </button>
        ) : null}
        {pendingAction.canAllIn ? (
          <button
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

      {pendingAction.canBet && isBetting ? (
        <div className="bet-panel">
          <label className="action-amount">
            <span>Bet amount</span>
            <input
              aria-label="Bet amount"
              type="number"
              min={getMinimumBet(pendingAction)}
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
              Confirm bet
            </button>
          </div>
        </div>
      ) : null}
    </section>
  );
}
