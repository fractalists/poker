type OrbitSpec = {
  minHeight: string;
  seatWidth: string;
  seatMinHeight: string;
  cardScale: number;
  boardTopPadding: string;
  boardSidePadding: string;
  boardBottomPadding: string;
  boardWidth: string;
  statWidth: string;
  boardGap: string;
};

export type OrbitPosition = {
  slot: string;
  x: number;
  y: number;
};

type SafeInset = {
  top: number;
  right: number;
  bottom: number;
  left: number;
};

export type OrbitLayout = {
  spec: OrbitSpec;
  positions: Map<number, OrbitPosition>;
  safeInset: SafeInset;
};

const baseTemplates: Record<number, OrbitPosition[]> = {
  2: [
    { slot: "hero", x: 50, y: 84 },
    { slot: "top-center", x: 50, y: 16 },
  ],
  3: [
    { slot: "hero", x: 50, y: 84 },
    { slot: "top-left", x: 36, y: 17 },
    { slot: "top-right", x: 64, y: 17 },
  ],
  4: [
    { slot: "hero", x: 50, y: 84 },
    { slot: "left-mid", x: 16, y: 46 },
    { slot: "top-center", x: 50, y: 16 },
    { slot: "right-mid", x: 84, y: 46 },
  ],
  5: [
    { slot: "hero", x: 50, y: 84 },
    { slot: "left-mid", x: 16, y: 46 },
    { slot: "top-left", x: 34, y: 18 },
    { slot: "top-right", x: 66, y: 18 },
    { slot: "right-mid", x: 84, y: 46 },
  ],
  6: [
    { slot: "hero", x: 50, y: 84 },
    { slot: "bottom-left", x: 28, y: 70 },
    { slot: "left-upper", x: 14, y: 38 },
    { slot: "top-center", x: 50, y: 13 },
    { slot: "right-upper", x: 86, y: 38 },
    { slot: "bottom-right", x: 72, y: 70 },
  ],
  7: [
    { slot: "hero", x: 50, y: 84 },
    { slot: "bottom-left", x: 32, y: 78 },
    { slot: "left-mid", x: 16, y: 46 },
    { slot: "top-left", x: 34, y: 17 },
    { slot: "top-right", x: 66, y: 17 },
    { slot: "right-mid", x: 84, y: 46 },
    { slot: "bottom-right", x: 68, y: 78 },
  ],
  8: [
    { slot: "hero", x: 50, y: 82 },
    { slot: "bottom-left", x: 34, y: 79 },
    { slot: "left-lower", x: 14, y: 58 },
    { slot: "left-upper", x: 14, y: 31 },
    { slot: "top-left", x: 36, y: 12 },
    { slot: "top-right", x: 64, y: 12 },
    { slot: "right-upper", x: 86, y: 31 },
    { slot: "right-lower", x: 86, y: 58 },
  ],
  9: [
    { slot: "hero", x: 50, y: 82 },
    { slot: "bottom-left", x: 33, y: 80 },
    { slot: "left-lower", x: 13, y: 60 },
    { slot: "left-upper", x: 13, y: 31 },
    { slot: "top-left", x: 34, y: 13 },
    { slot: "top-center", x: 50, y: 9 },
    { slot: "top-right", x: 66, y: 13 },
    { slot: "right-upper", x: 87, y: 31 },
    { slot: "right-lower", x: 87, y: 60 },
  ],
  10: [
    { slot: "hero", x: 50, y: 82 },
    { slot: "bottom-left", x: 33, y: 80 },
    { slot: "left-lower", x: 13, y: 60 },
    { slot: "left-upper", x: 13, y: 31 },
    { slot: "top-left", x: 34, y: 13 },
    { slot: "top-center", x: 50, y: 9 },
    { slot: "top-right", x: 66, y: 13 },
    { slot: "right-upper", x: 87, y: 31 },
    { slot: "right-lower", x: 87, y: 60 },
    { slot: "bottom-right", x: 67, y: 80 },
  ],
};

function clamp(value: number, min: number, max: number) {
  return Math.min(Math.max(value, min), max);
}

function roundPercent(value: number) {
  return Math.round(value * 100) / 100;
}

function toPx(value: number) {
  return `${Math.round(value)}px`;
}

function buildOrbitSpec(playerCount: number): OrbitSpec {
  const normalized = clamp(playerCount, 2, 10);
  const crowding = (normalized - 2) / 8;

  const minHeight = 560 + crowding * 244;
  const seatWidth = 208 - crowding * 24;
  const seatMinHeight = 160 - crowding * 8;
  const cardScale = 1 - crowding * 0.05;
  const balancedBoardPadding = 96 + crowding * 80;
  const crowdedCueBias = Math.max(0, crowding - 0.5) * 114;
  const boardTopPadding = balancedBoardPadding + crowdedCueBias;
  const boardSidePadding = 72 + crowding * 94;
  const boardBottomPadding = balancedBoardPadding - crowdedCueBias;
  const boardWidth = 700 - crowding * 56;
  const statWidth = 132 - crowding * 18;
  const boardGap = 24 - crowding * 2;

  return {
    minHeight: toPx(minHeight),
    seatWidth: toPx(seatWidth),
    seatMinHeight: toPx(seatMinHeight),
    cardScale: Number(cardScale.toFixed(2)),
    boardTopPadding: toPx(boardTopPadding),
    boardSidePadding: toPx(boardSidePadding),
    boardBottomPadding: toPx(boardBottomPadding),
    boardWidth: toPx(boardWidth),
    statWidth: toPx(statWidth),
    boardGap: toPx(boardGap),
  };
}

function px(value: string) {
  return Number.parseFloat(value.replace("px", ""));
}

function buildSafeInset(spec: OrbitSpec): SafeInset {
  const minHeight = px(spec.minHeight);
  const seatWidth = px(spec.seatWidth);
  const seatHeight = px(spec.seatMinHeight);
  const boardWidth = px(spec.boardWidth);
  const boardSidePadding = px(spec.boardSidePadding);

  const estimatedStageWidth = boardWidth + boardSidePadding * 2;
  const verticalHalfSeatPercent = (seatHeight / minHeight) * 50;
  const horizontalHalfSeatPercent = (seatWidth / estimatedStageWidth) * 50;
  const isCompactTable = minHeight < 720;
  const topClearanceOffset = isCompactTable ? 2.5 : 1.6;
  const bottomClearanceOffset = isCompactTable ? 0.5 : 2.5;

  return {
    top: roundPercent(
      clamp(verticalHalfSeatPercent + topClearanceOffset, 10.25, 17),
    ),
    right: roundPercent(clamp(horizontalHalfSeatPercent + 2.5, 10.5, 15)),
    bottom: roundPercent(
      clamp(verticalHalfSeatPercent + bottomClearanceOffset, 10.5, 18),
    ),
    left: roundPercent(clamp(horizontalHalfSeatPercent + 2.5, 10.5, 15)),
  };
}

function fitTemplateToSafeInset(
  template: OrbitPosition[],
  safeInset: SafeInset,
): OrbitPosition[] {
  const minX = Math.min(...template.map((slot) => slot.x));
  const maxX = Math.max(...template.map((slot) => slot.x));
  const minY = Math.min(...template.map((slot) => slot.y));
  const maxY = Math.max(...template.map((slot) => slot.y));
  const innerWidth = 100 - safeInset.left - safeInset.right;
  const innerHeight = 100 - safeInset.top - safeInset.bottom;

  return template.map((slot) => ({
    slot: slot.slot,
    x: roundPercent(
      safeInset.left + ((slot.x - minX) / (maxX - minX || 1)) * innerWidth,
    ),
    y: roundPercent(
      safeInset.top + ((slot.y - minY) / (maxY - minY || 1)) * innerHeight,
    ),
  }));
}

function assignSeatsToPositions(
  playerCount: number,
  humanSeat: number | undefined,
  template: OrbitPosition[],
) {
  const positions = new Map<number, OrbitPosition>();
  const normalizedPlayerCount = clamp(playerCount, 2, 10);
  const normalizedHumanSeat =
    humanSeat !== undefined &&
    humanSeat >= 0 &&
    humanSeat < normalizedPlayerCount
      ? humanSeat
      : normalizedPlayerCount - 1;

  template.forEach((slot, offset) => {
    const seatIndex =
      offset === 0
        ? normalizedHumanSeat
        : (normalizedHumanSeat + offset) % normalizedPlayerCount;
    positions.set(seatIndex, slot);
  });

  return positions;
}

export function buildOrbitLayout(
  playerCount: number,
  humanSeat?: number,
): OrbitLayout {
  const normalizedPlayerCount = clamp(playerCount, 2, 10);
  const spec = buildOrbitSpec(normalizedPlayerCount);
  const safeInset = buildSafeInset(spec);
  const baseTemplate =
    baseTemplates[normalizedPlayerCount] ?? baseTemplates[normalizedPlayerCount - 1];
  const fittedTemplate = fitTemplateToSafeInset(baseTemplate, safeInset);
  const positions = assignSeatsToPositions(
    normalizedPlayerCount,
    humanSeat,
    fittedTemplate,
  );

  return {
    spec,
    positions,
    safeInset,
  };
}
