import { render, screen } from "@testing-library/react";

import { TableSeat } from "./TableSeat";

describe("TableSeat", () => {
  it("adds suit-specific classes for visible hole cards", () => {
    render(
      <TableSeat
        seat={{
          index: 5,
          name: "Player6",
          status: "PLAYING",
          bankroll: 100,
          inPotAmount: 2,
          isTurn: false,
          cards: ["♥A", "♠10"],
        }}
        viewerSeat={5}
      />,
    );

    const heartCard = screen
      .getAllByText((_, element) => element?.textContent === "♥A")
      .find((element) => element.classList.contains("card-face"));
    const spadeCard = screen
      .getAllByText((_, element) => element?.textContent === "♠10")
      .find((element) => element.classList.contains("card-face"));

    expect(heartCard).toHaveClass("card-face", "card-face--hearts");
    expect(spadeCard).toHaveClass("card-face", "card-face--spades");
  });

  it("renders bankroll and settlement delta in the same footer group", () => {
    const { container } = render(
      <TableSeat
        seat={{
          index: 0,
          name: "Player1",
          status: "PLAYING",
          bankroll: 118,
          inPotAmount: 0,
          isTurn: false,
          cards: ["♥A", "♠K"],
          netChange: 18,
          bestHand: "Straight",
          isWinner: true,
        }}
      />,
    );

    const bankrollGroup = container.querySelector(".seat-bankroll-group");
    const bankrollPill = container.querySelector(".seat-stack-pill--bankroll");

    expect(bankrollGroup).toBeInTheDocument();
    expect(bankrollPill).toHaveTextContent("Bankroll");
    expect(bankrollPill).toHaveTextContent("118");
    expect(bankrollGroup).toHaveTextContent("+18");
    expect(bankrollPill).toBeInTheDocument();
    expect(bankrollPill?.querySelector(".seat-result-pill")).toBeInTheDocument();
    expect(bankrollGroup?.querySelector(":scope > .seat-result-pill")).not.toBeInTheDocument();
    expect(container.querySelector(".seat-body")).toBeInTheDocument();
    expect(container.querySelector(".seat-footer")).toBeInTheDocument();
  });

  it("marks the viewer's own seat inline in the seat header", () => {
    render(
      <TableSeat
        seat={{
          index: 5,
          name: "Player6",
          status: "PLAYING",
          bankroll: 100,
          inPotAmount: 2,
          isTurn: false,
          cards: ["♥A", "♠10"],
        }}
        viewerSeat={5}
      />,
    );

    expect(screen.getByText("Player6")).toBeInTheDocument();
    expect(screen.getByText("You")).toBeInTheDocument();
    expect(screen.queryByText("Seat 6")).not.toBeInTheDocument();
    expect(screen.queryByText("PLAYING")).not.toBeInTheDocument();
  });

  it("shows the standard table position abbreviation in the seat header", () => {
    const { container } = render(
      <TableSeat
        seat={{
          index: 0,
          name: "Player1",
          position: "UTG",
          status: "PLAYING",
          bankroll: 100,
          inPotAmount: 2,
          isTurn: false,
          cards: ["**", "**"],
        }}
      />,
    );

    expect(screen.getByText("UTG")).toBeInTheDocument();
    expect(container.querySelector(".seat-position-badge")).toHaveTextContent("UTG");
    expect(screen.queryByText("Seat 1")).not.toBeInTheDocument();
  });

  it("renders folded seats dimmed and current-turn seats highlighted", () => {
    const { container, rerender } = render(
      <TableSeat
        seat={{
          index: 3,
          name: "Player4",
          status: "OUT",
          bankroll: 100,
          inPotAmount: 0,
          isTurn: false,
          cards: ["**", "**"],
        }}
        recentAction={{ label: "Folded", tone: "fold", stamp: "seat-3-fold" }}
      />,
    );

    expect(container.querySelector(".table-seat.is-folded")).toBeInTheDocument();
    expect(screen.getByText("Folded")).toBeInTheDocument();
    expect(
      container.querySelector(".seat-header-meta .seat-action-pill.tone-fold"),
    ).toBeInTheDocument();
    expect(
      container.querySelector(".seat-meta-row .seat-action-pill"),
    ).not.toBeInTheDocument();
    expect(screen.queryByText("OUT")).not.toBeInTheDocument();
    expect(screen.queryByText("PLAYING")).not.toBeInTheDocument();

    rerender(
      <TableSeat
        seat={{
          index: 5,
          name: "Player6",
          status: "PLAYING",
          bankroll: 89,
          inPotAmount: 0,
          isTurn: true,
          cards: ["♠J", "♦7"],
        }}
      />,
    );

    expect(container.querySelector(".table-seat.is-turn")).toBeInTheDocument();
    expect(screen.getByText("To act")).toBeInTheDocument();
  });

  it("shows all-in seats with a dedicated table badge", () => {
    const { container } = render(
      <TableSeat
        seat={{
          index: 2,
          name: "Player3",
          status: "ALLIN",
          bankroll: 0,
          inPotAmount: 32,
          isTurn: false,
          cards: ["**", "**"],
        }}
      />,
    );

    expect(container.querySelector(".table-seat.is-all-in")).toBeInTheDocument();
    expect(screen.getByText("All-in")).toBeInTheDocument();
  });

  it("renders an elimination badge for busted seats even without a fresh action", () => {
    const { container } = render(
      <TableSeat
        seat={{
          index: 8,
          name: "Player9",
          status: "OUT",
          bankroll: 0,
          inPotAmount: 0,
          isTurn: false,
          cards: ["**", "**"],
        }}
      />,
    );

    expect(container.querySelector(".table-seat.is-eliminated")).toBeInTheDocument();
    expect(container.querySelector(".table-seat.is-folded")).not.toBeInTheDocument();
    expect(
      container.querySelector(".seat-header-meta .seat-action-pill.tone-out"),
    ).toBeInTheDocument();
    expect(screen.getByText("Busted")).toBeInTheDocument();
  });

  it("adds dedicated settlement emphasis classes for winners and losers", () => {
    const { container, rerender } = render(
      <TableSeat
        seat={{
          index: 0,
          name: "Player1",
          status: "PLAYING",
          bankroll: 124,
          inPotAmount: 0,
          isTurn: false,
          cards: ["♥A", "♠K"],
          netChange: 24,
          bestHand: "Straight",
          isWinner: true,
        }}
        showSettlementEffects
      />,
    );

    expect(
      container.querySelector(".table-seat.is-settlement-winner"),
    ).toBeInTheDocument();
    expect(
      container.querySelector(".table-seat.is-settlement-loser"),
    ).not.toBeInTheDocument();

    rerender(
      <TableSeat
        seat={{
          index: 5,
          name: "Player6",
          status: "PLAYING",
          bankroll: 76,
          inPotAmount: 0,
          isTurn: false,
          cards: ["♣9", "♦9"],
          netChange: -24,
          bestHand: "One pair",
        }}
        showSettlementEffects
      />,
    );

    expect(
      container.querySelector(".table-seat.is-settlement-loser"),
    ).toBeInTheDocument();
    expect(
      container.querySelector(".table-seat.is-settlement-winner"),
    ).not.toBeInTheDocument();
  });
});
