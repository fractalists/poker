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
    expect(screen.getByText("Seat 6")).toBeInTheDocument();
    expect(screen.queryByText("PLAYING")).not.toBeInTheDocument();
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
});
