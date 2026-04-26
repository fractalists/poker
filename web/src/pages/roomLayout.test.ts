import { buildOrbitLayout } from "./roomLayout";

function px(value: string) {
  return Number.parseFloat(value.replace("px", ""));
}

function verticalHalfSeatPercent(layout: ReturnType<typeof buildOrbitLayout>) {
  const renderedSeatOverflow = 32;
  return (
    ((px(layout.spec.seatMinHeight) + renderedSeatOverflow) /
      px(layout.spec.minHeight)) *
    50
  );
}

function boardLaneCenterPercent(layout: ReturnType<typeof buildOrbitLayout>) {
  const minHeight = px(layout.spec.minHeight);
  const topPadding = px(layout.spec.boardTopPadding);
  const bottomPadding = px(layout.spec.boardBottomPadding);

  return (
    ((topPadding + (minHeight - topPadding - bottomPadding) / 2) /
      minHeight) *
    100
  );
}

describe("roomLayout", () => {
  it("keeps seat anchors inside the computed safe inset for 2, 6, and 10 handed tables", () => {
    [2, 6, 10].forEach((playerCount) => {
      const layout = buildOrbitLayout(playerCount, playerCount - 1);

      expect(layout.positions.size).toBe(playerCount);

      for (const position of layout.positions.values()) {
        expect(position.x).toBeGreaterThanOrEqual(layout.safeInset.left);
        expect(position.x).toBeLessThanOrEqual(100 - layout.safeInset.right);
        expect(position.y).toBeGreaterThanOrEqual(layout.safeInset.top);
        expect(position.y).toBeLessThanOrEqual(100 - layout.safeInset.bottom);
      }
    });
  });

  it("keeps full-ring seats above the readability floor", () => {
    const layout = buildOrbitLayout(10, 9);

    expect(px(layout.spec.seatWidth)).toBeGreaterThanOrEqual(184);
    expect(px(layout.spec.seatMinHeight)).toBeGreaterThanOrEqual(152);
    expect(layout.spec.cardScale).toBeGreaterThanOrEqual(0.94);
    expect(px(layout.spec.minHeight)).toBeGreaterThanOrEqual(804);
  });

  it("packs six-handed tables into a denser but still readable seat footprint", () => {
    const layout = buildOrbitLayout(6, 5);

    expect(px(layout.spec.seatWidth)).toBeLessThanOrEqual(200);
    expect(px(layout.spec.seatMinHeight)).toBeLessThanOrEqual(158);
    expect(px(layout.spec.minHeight)).toBeLessThanOrEqual(684);
  });

  it("keeps six-handed top and bottom seats clear of the center board HUD lane", () => {
    const layout = buildOrbitLayout(6, 5);
    const seatHalf = verticalHalfSeatPercent(layout);
    const topSeat = [...layout.positions.values()].find(
      (position) => position.slot === "top-center",
    );
    const heroSeat = [...layout.positions.values()].find(
      (position) => position.slot === "hero",
    );

    expect(topSeat).toBeDefined();
    expect(heroSeat).toBeDefined();
    expect(topSeat!.y + seatHalf).toBeLessThanOrEqual(24.6);
    expect(heroSeat!.y - seatHalf).toBeGreaterThanOrEqual(74);
  });

  it("keeps full-ring community cards visually centered between top and bottom seats", () => {
    const layout = buildOrbitLayout(10, 9);

    expect(px(layout.spec.boardTopPadding)).toBeGreaterThan(
      px(layout.spec.boardBottomPadding),
    );
    expect(boardLaneCenterPercent(layout)).toBeGreaterThanOrEqual(56.5);
    expect(boardLaneCenterPercent(layout)).toBeLessThanOrEqual(57.5);
  });
});
