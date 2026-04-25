import { useEffect, useRef, useState } from "react";

import { CardFace } from "./CardFace";

export const COMMUNITY_CARD_REVEAL_INTERVAL_MS = 220;
export const COMMUNITY_CARD_FLIP_MS = 420;

const faceDownTokens = new Set(["**", "--"]);

type CommunityCardsProps = {
  cards: string[];
};

function isFaceDown(card: string | undefined) {
  return card === undefined || faceDownTokens.has(card) || card.trim() === "";
}

function shouldRevealCard(previous: string | undefined, next: string) {
  return !isFaceDown(next) && (isFaceDown(previous) || previous !== next);
}

function clearTimers(timers: number[]) {
  for (const timer of timers) {
    window.clearTimeout(timer);
  }
  timers.length = 0;
}

export function CommunityCards({ cards }: CommunityCardsProps) {
  const previousCardsRef = useRef(cards);
  const timersRef = useRef<number[]>([]);
  const [displayCards, setDisplayCards] = useState(cards);
  const [revealingIndexes, setRevealingIndexes] = useState<Set<number>>(
    () => new Set(),
  );
  const cardsKey = cards.join("|");

  useEffect(() => {
    const previousCards = previousCardsRef.current;
    const revealIndexes = cards
      .map((card, index) => ({
        card,
        index,
        shouldReveal: shouldRevealCard(previousCards[index], card),
      }))
      .filter((entry) => entry.shouldReveal);

    clearTimers(timersRef.current);

    if (revealIndexes.length === 0) {
      previousCardsRef.current = cards;
      setDisplayCards(cards);
      setRevealingIndexes(new Set());
      return;
    }

    previousCardsRef.current = cards;
    setDisplayCards(
      cards.map((card, index) =>
        revealIndexes.some((entry) => entry.index === index) ? "**" : card,
      ),
    );
    setRevealingIndexes(new Set());

    revealIndexes.forEach((entry, order) => {
      const revealTimer = window.setTimeout(() => {
        setDisplayCards((currentCards) =>
          currentCards.map((card, index) =>
            index === entry.index ? entry.card : card,
          ),
        );
        setRevealingIndexes((currentIndexes) => {
          const nextIndexes = new Set(currentIndexes);
          nextIndexes.add(entry.index);
          return nextIndexes;
        });

        const settleTimer = window.setTimeout(() => {
          setRevealingIndexes((currentIndexes) => {
            const nextIndexes = new Set(currentIndexes);
            nextIndexes.delete(entry.index);
            return nextIndexes;
          });
        }, COMMUNITY_CARD_FLIP_MS);
        timersRef.current.push(settleTimer);
      }, COMMUNITY_CARD_REVEAL_INTERVAL_MS * (order + 1));

      timersRef.current.push(revealTimer);
    });

    return () => clearTimers(timersRef.current);
  }, [cardsKey]);

  return (
    <>
      {displayCards.map((card, index) => (
        <div
          className={[
            "board-card",
            "community-card",
            revealingIndexes.has(index) ? "is-revealing" : "",
          ]
            .filter(Boolean)
            .join(" ")}
          key={`community-${index}`}
        >
          <CardFace card={card} />
        </div>
      ))}
    </>
  );
}
