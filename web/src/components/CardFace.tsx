const suitClassMap: Record<string, string> = {
  "♥": "card-face--hearts",
  "♦": "card-face--diamonds",
  "♣": "card-face--clubs",
  "♠": "card-face--spades",
};

const faceDownTokens = new Set(["**", "--"]);

type CardFaceProps = {
  card: string;
};

export function CardFace({ card }: CardFaceProps) {
  if (faceDownTokens.has(card) || card.trim() === "") {
    return (
      <span aria-label="face-down card" className="card-face card-face--back">
        <span aria-hidden="true" className="card-face-pattern" />
      </span>
    );
  }

  const suit = card.charAt(0);
  const rank = card.slice(1);
  const suitClass = suitClassMap[suit];

  if (!suitClass || rank.length === 0) {
    return <span className="card-face card-face--unknown">{card}</span>;
  }

  return (
    <span className={`card-face ${suitClass}`}>
      <span className="card-face-suit">{suit}</span>
      <span className="card-face-rank">{rank}</span>
    </span>
  );
}
