import {
  act,
  fireEvent,
  render,
  screen,
  waitFor,
  within,
} from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { vi } from "vitest";

import type { RoomSnapshot } from "../lib/types";
import { RoomPage, RoomRoute } from "./RoomPage";
import { ACTION_PLAYBACK_DELAY_MS } from "./roomPlayback";

class MockWebSocket {
  static instances: MockWebSocket[] = [];

  readonly close = vi.fn();
  onclose: (() => void) | null = null;
  onerror: (() => void) | null = null;
  onopen: (() => void) | null = null;
  onmessage: ((event: MessageEvent) => void) | null = null;
  readonly url: string;

  constructor(url: string | URL) {
    this.url = String(url);
    MockWebSocket.instances.push(this);
  }
}

describe("RoomPage", () => {
  afterEach(() => {
    MockWebSocket.instances = [];
    vi.unstubAllGlobals();
    window.localStorage.clear();
  });

  it("renders the live table, board cards, and seat states", () => {
    const { container } = render(
      <RoomPage
        snapshot={{
          roomId: "room-001",
          roomName: "Table 1",
          status: "awaiting_action",
          viewerRole: "player",
          handNumber: 3,
          smallBlind: 1,
          pot: 6,
          currentAmount: 2,
          round: "FLOP",
          boardCards: ["♥Q", "**", "**"],
          seats: [
            {
              index: 0,
              name: "Player1",
              status: "PLAYING",
              bankroll: 98,
              inPotAmount: 2,
              isTurn: false,
              cards: ["**", "**"],
            },
            {
              index: 5,
              name: "Player6",
              status: "PLAYING",
              bankroll: 99,
              inPotAmount: 1,
              isTurn: true,
              cards: ["♣A", "♣K"],
            },
          ],
          pendingAction: {
            token: "turn-1",
            seatIndex: 5,
            minAmount: 1,
            maxAmount: 99,
            canCheck: false,
            canCall: true,
            canBet: true,
            canFold: true,
            canAllIn: true,
          },
        }}
        onAction={async () => {}}
        onStartHand={async () => {}}
        onTakeSeat={async () => {}}
      />,
    );

    const flopCard = screen
      .getAllByText((_, element) => element?.textContent === "♥Q")
      .find((element) => element.classList.contains("card-face"));

    expect(screen.getByText("Table 1")).toBeInTheDocument();
    expect(flopCard).toBeInTheDocument();
    expect(screen.getByText("Player6")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /fold/i })).toBeInTheDocument();
    expect(flopCard).toHaveClass("card-face", "card-face--hearts");
    expect(container.querySelector(".room-shell--fixed")).toBeInTheDocument();
    expect(
      container.querySelector(".table-stage > .room-stats"),
    ).not.toBeInTheDocument();
    expect(container.querySelector(".board-cluster")).toBeInTheDocument();
    const boardCenterpiece = container.querySelector(".board-centerpiece");
    const boardCardsShell = container.querySelector(".board-cards-shell");
    const boardMeta = container.querySelector(".board-meta");
    expect(boardMeta).toBeInTheDocument();
    expect(boardMeta?.parentElement).toBe(boardCenterpiece);
    expect(
      Array.from(boardCenterpiece?.children ?? []).map(
        (element) => element.className,
      ),
    ).toEqual(["board-badge", "board-cards-shell", "board-meta"]);
    expect(boardCardsShell?.nextElementSibling).toBe(boardMeta);
    expect(container.querySelectorAll(".table-stat-row")).toHaveLength(3);
    expect(container.querySelector(".table-stat-row--pot")).toHaveTextContent("6");
    expect(container.querySelector(".table-stat-row--current")).toHaveTextContent("2");
    expect(container.querySelector(".table-stat-icon")).not.toBeInTheDocument();
    expect(container.querySelector(".board-badge-dot")).not.toBeInTheDocument();
    expect(container.querySelector(".seat-orbit")).toBeInTheDocument();
    expect(container.querySelector(".seat-grid")).not.toBeInTheDocument();
    expect(container.querySelector(".side-panel--scroll")).toBeInTheDocument();
    expect(
      container.querySelector(".room-history-stack.room-history-stack--scroll"),
    ).toBeInTheDocument();
    expect(screen.queryByText("PLAYING")).not.toBeInTheDocument();
  });

  it("keeps the player's hand and call context inside the action controls", () => {
    const { container } = render(
      <RoomPage
        snapshot={{
          roomId: "room-001",
          roomName: "Decision Table",
          status: "awaiting_action",
          viewerRole: "player",
          humanSeat: 5,
          playerCount: 6,
          handNumber: 3,
          smallBlind: 1,
          pot: 94,
          currentAmount: 83,
          round: "PREFLOP",
          boardCards: ["**", "**", "**", "**", "**"],
          seats: [
            {
              index: 5,
              name: "Player6",
              status: "PLAYING",
              bankroll: 100,
              inPotAmount: 0,
              isTurn: true,
              cards: ["♣A", "♠2"],
            },
          ],
          pendingAction: {
            token: "turn-1",
            seatIndex: 5,
            minAmount: 83,
            minBetAmount: 100,
            maxAmount: 100,
            canCheck: false,
            canCall: true,
            canBet: true,
            canFold: true,
            canAllIn: true,
          },
        }}
        onAction={async () => {}}
        onStartHand={async () => {}}
        onTakeSeat={async () => {}}
      />,
    );

    const actionBar = container.querySelector(".action-bar") as HTMLElement;
    const actionBarQueries = within(actionBar);

    expect(actionBarQueries.getByText("Your hand")).toBeInTheDocument();
    expect(
      actionBar.querySelectorAll(".decision-card .card-face"),
    ).toHaveLength(2);
    expect(actionBarQueries.queryByText("To call")).not.toBeInTheDocument();
    expect(actionBarQueries.getByText("Stack")).toBeInTheDocument();
    expect(actionBarQueries.getByText("100")).toBeInTheDocument();
    expect(actionBarQueries.getByText("After call")).toBeInTheDocument();
    expect(actionBarQueries.getByText("17")).toBeInTheDocument();
    expect(actionBarQueries.getByText("Pot odds")).toBeInTheDocument();
    expect(actionBarQueries.getByText("47%")).toBeInTheDocument();
    expect(actionBarQueries.getByRole("button", { name: "Call 83" })).toBeInTheDocument();
  });

  it("uses a full-ring seat orbit when the room has ten players", () => {
    const { container } = render(
      <RoomPage
        snapshot={{
          roomId: "room-010",
          roomName: "Full Ring",
          status: "running",
          viewerRole: "player",
          humanSeat: 9,
          playerCount: 10,
          handNumber: 2,
          smallBlind: 1,
          pot: 12,
          currentAmount: 4,
          round: "PREFLOP",
          boardCards: [],
          seats: Array.from({ length: 10 }, (_, index) => ({
            index,
            name: `Player${index + 1}`,
            status: "PLAYING",
            bankroll: 100 - index,
            inPotAmount: index === 0 ? 1 : 0,
            isTurn: index === 9,
            cards: index === 9 ? ["♠A", "♥K"] : ["**", "**"],
          })),
        }}
        onAction={async () => {}}
        onStartHand={async () => {}}
        onTakeSeat={async () => {}}
      />,
    );

    expect(container.querySelector(".table-live-layout")).toHaveStyle({
      "--orbit-min-height": "804px",
      "--orbit-seat-width": "184px",
    });
    expect(container.querySelector(".seat-orbit--10")).toBeInTheDocument();
    expect(container.querySelectorAll(".seat-slot")).toHaveLength(10);
    expect(container.querySelector(".seat-slot--hero")).toBeInTheDocument();
    expect(
      parseFloat(
        (
          container.querySelector(".seat-slot--top-center") as HTMLElement
        ).style.top,
      ),
    ).toBeGreaterThanOrEqual(11);
    expect(
      parseFloat(
        (
          container.querySelector(".seat-slot--top-left") as HTMLElement
        ).style.top,
      ),
    ).toBeGreaterThanOrEqual(14);
    expect(
      parseFloat(
        (
          container.querySelector(".seat-slot--left-lower") as HTMLElement
        ).style.left,
      ),
    ).toBeGreaterThanOrEqual(11.5);
    expect(
      parseFloat(
        (
          container.querySelector(".seat-slot--bottom-left") as HTMLElement
        ).style.top,
      ),
    ).toBeGreaterThanOrEqual(78);
    expect(
      parseFloat(
        (
          container.querySelector(".seat-slot--right-lower") as HTMLElement
        ).style.left,
      ),
    ).toBeLessThanOrEqual(88.5);
  });

  it("prioritizes the player's seat and current actor in the narrow-screen seat order", () => {
    render(
      <RoomPage
        snapshot={{
          roomId: "room-010",
          roomName: "Full Ring",
          status: "awaiting_action",
          viewerRole: "player",
          humanSeat: 9,
          playerCount: 10,
          handNumber: 8,
          smallBlind: 1,
          pot: 20,
          currentAmount: 4,
          round: "TURN",
          boardCards: ["♠A", "♥K", "♣4", "♦9", "**"],
          seats: Array.from({ length: 10 }, (_, index) => ({
            index,
            name: `Player${index + 1}`,
            status:
              index === 7 || index === 8
                ? "OUT"
                : "PLAYING",
            bankroll: index === 8 ? 0 : 100 - index,
            inPotAmount: index === 3 ? 4 : 0,
            isTurn: index === 3,
            cards: index === 9 ? ["♠A", "♥K"] : ["**", "**"],
          })),
          pendingAction: {
            token: "turn-8",
            seatIndex: 3,
            minAmount: 4,
            maxAmount: 96,
            canCheck: false,
            canCall: true,
            canBet: true,
            canFold: true,
            canAllIn: true,
          },
        }}
        onAction={async () => {}}
        onStartHand={async () => {}}
        onTakeSeat={async () => {}}
      />,
    );

    const ownSlot = screen.getByText("Player10").closest(".seat-slot") as HTMLElement;
    const currentSlot = screen.getByText("Player4").closest(".seat-slot") as HTMLElement;
    const foldedSlot = screen.getByText("Player8").closest(".seat-slot") as HTMLElement;
    const bustedSlot = screen.getByText("Player9").closest(".seat-slot") as HTMLElement;

    expect(ownSlot.style.getPropertyValue("--mobile-seat-order")).toBe("0");
    expect(currentSlot.style.getPropertyValue("--mobile-seat-order")).toBe("1");
    expect(Number(foldedSlot.style.getPropertyValue("--mobile-seat-order"))).toBeLessThan(
      Number(bustedSlot.style.getPropertyValue("--mobile-seat-order")),
    );
  });

  it("renders waiting rooms without crashing when collections are missing", () => {
    const { container } = render(
      <RoomPage
        snapshot={
          {
            roomId: "room-001",
            roomName: "Table 1",
            status: "waiting",
            handNumber: 0,
            smallBlind: 1,
          } as RoomSnapshot
        }
        onAction={async () => {}}
        onStartHand={async () => {}}
        onTakeSeat={async () => {}}
      />,
    );

    expect(screen.getByText("Table 1")).toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: /start hand/i }),
    ).toBeInTheDocument();
    expect(screen.getByText(/no events yet/i)).toBeInTheDocument();
    expect(
      screen.getByText(/waiting for the opening deal/i),
    ).toBeInTheDocument();
    expect(
      screen.getByText(/start a hand to bring players onto the felt/i),
    ).toBeInTheDocument();
    expect(container.querySelector(".table-empty-state")).toBeInTheDocument();
    expect(container.querySelector(".board-cluster")).not.toBeInTheDocument();
    expect(
      container.querySelector(".table-placeholder"),
    ).not.toBeInTheDocument();
  });

  it("renders seats whose cards have not been dealt yet", () => {
    render(
      <RoomPage
        snapshot={{
          roomId: "room-001",
          roomName: "Table 1",
          status: "running",
          viewerRole: "player",
          humanSeat: 5,
          handNumber: 1,
          smallBlind: 1,
          pot: 3,
          currentAmount: 2,
          round: "PREFLOP",
          boardCards: [],
          seats: [
            {
              index: 5,
              name: "Player6",
              status: "PLAYING",
              bankroll: 99,
              inPotAmount: 1,
              isTurn: false,
              cards: undefined as unknown as string[],
            },
          ],
        }}
        onAction={async () => {}}
        onStartHand={async () => {}}
        onTakeSeat={async () => {}}
      />,
    );

    expect(screen.getByText("Player6")).toBeInTheDocument();
    expect(screen.getByText("You")).toBeInTheDocument();
    expect(screen.queryByText(/player view active/i)).not.toBeInTheDocument();
    expect(
      screen.queryByText(/your hole cards are shown/i),
    ).not.toBeInTheDocument();
  });

  it("renders a hand-finished summary with winners and chip swings", () => {
    const { container } = render(
      <RoomPage
        snapshot={{
          roomId: "room-001",
          roomName: "Table 1",
          status: "hand_finished",
          viewerRole: "player",
          handNumber: 4,
          humanSeat: 5,
          smallBlind: 1,
          pot: 0,
          currentAmount: 0,
          round: "FINISH",
          boardCards: ["♥Q", "♣J", "♦10", "♠2", "♣3"],
          seats: [
            {
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
            },
            {
              index: 2,
              name: "Player3",
              status: "PLAYING",
              bankroll: 100,
              inPotAmount: 0,
              isTurn: false,
              cards: ["**", "**"],
              netChange: 0,
              bestHand: "No pair",
              isWinner: false,
            },
            {
              index: 5,
              name: "Player6",
              status: "PLAYING",
              bankroll: 82,
              inPotAmount: 0,
              isTurn: false,
              cards: ["♣9", "♦9"],
              netChange: -18,
              bestHand: "One pair",
              isWinner: false,
            },
          ],
        }}
        onAction={async () => {}}
        onStartHand={async () => {}}
        onTakeSeat={async () => {}}
      />,
    );

    expect(screen.getByText("Player1 wins +18")).toBeInTheDocument();
    expect(screen.getByText("Winning hand: Straight")).toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: /start next hand/i }),
    ).toBeInTheDocument();
    expect(screen.getAllByText("Player6")).toHaveLength(2);
    expect(screen.getAllByText("-18")).toHaveLength(2);
    expect(screen.getAllByText("One pair")).toHaveLength(2);
    const historyStack = container.querySelector(
      ".room-history-stack.room-history-stack--scroll",
    );
    const settlementPanel = container.querySelector(
      ".settlement-panel.is-animated",
    );
    const feedPanel = container.querySelector(".room-feed-panel");
    const settlementQueries = within(settlementPanel as HTMLElement);
    expect(historyStack).toBeInTheDocument();
    expect(settlementPanel?.parentElement).toBe(historyStack);
    expect(feedPanel?.parentElement).toBe(historyStack);
    expect(
      container.querySelector(".table-seat.is-settlement-winner"),
    ).toBeInTheDocument();
    expect(
      container.querySelector(".table-seat.is-settlement-loser"),
    ).toBeInTheDocument();
    expect(
      container.querySelectorAll(".settlement-entry").length,
    ).toBe(2);
    expect(settlementQueries.queryByText("Player3")).not.toBeInTheDocument();
    expect(settlementQueries.queryByText("No pair")).not.toBeInTheDocument();
    expect(settlementQueries.queryByText(/^0$/)).not.toBeInTheDocument();
  });

  it("reveals newly dealt community cards one at a time", async () => {
    const baseSnapshot: RoomSnapshot = {
      roomId: "room-001",
      roomName: "Table 1",
      status: "running",
      viewerRole: "player",
      humanSeat: 1,
      playerCount: 2,
      handNumber: 6,
      smallBlind: 1,
      pot: 3,
      currentAmount: 2,
      round: "PREFLOP",
      boardCards: ["**", "**", "**", "**", "**"],
      seats: [
        {
          index: 0,
          name: "Player1",
          status: "PLAYING",
          bankroll: 98,
          inPotAmount: 2,
          isTurn: false,
          cards: ["**", "**"],
        },
        {
          index: 1,
          name: "Player2",
          status: "PLAYING",
          bankroll: 99,
          inPotAmount: 1,
          isTurn: true,
          cards: ["♣A", "♣K"],
        },
      ],
    };

    const findBoardCard = (container: HTMLElement, card: string) =>
      Array.from(container.querySelectorAll(".board-card .card-face")).find(
        (element) => element.textContent === card,
      );

    const { container, rerender } = render(
      <RoomPage
        snapshot={baseSnapshot}
        onAction={async () => {}}
        onStartHand={async () => {}}
        onTakeSeat={async () => {}}
      />,
    );

    vi.useFakeTimers();
    try {
      rerender(
        <RoomPage
          snapshot={{
            ...baseSnapshot,
            round: "FLOP",
            boardCards: ["♥Q", "♣J", "♦10", "**", "**"],
          }}
          onAction={async () => {}}
          onStartHand={async () => {}}
          onTakeSeat={async () => {}}
        />,
      );

      const board = container.querySelector(".board-cards") as HTMLElement;

      expect(findBoardCard(board, "♥Q")).toBeUndefined();
      expect(board.querySelectorAll(".card-face--back")).toHaveLength(5);

      await act(async () => {
        await vi.advanceTimersByTimeAsync(220);
      });

      expect(findBoardCard(board, "♥Q")).toBeInTheDocument();
      expect(findBoardCard(board, "♣J")).toBeUndefined();

      await act(async () => {
        await vi.advanceTimersByTimeAsync(220);
      });

      expect(findBoardCard(board, "♣J")).toBeInTheDocument();
      expect(findBoardCard(board, "♦10")).toBeUndefined();

      await act(async () => {
        await vi.advanceTimersByTimeAsync(220);
      });

      expect(findBoardCard(board, "♦10")).toBeInTheDocument();
    } finally {
      vi.useRealTimers();
    }
  });

  it("requires an inline confirmation before leaving for the lobby", () => {
    render(
      <RoomPage
        snapshot={{
          roomId: "room-001",
          roomName: "Table 1",
          status: "awaiting_action",
          viewerRole: "player",
          handNumber: 3,
          humanSeat: 5,
          smallBlind: 1,
          pot: 6,
          currentAmount: 2,
          round: "FLOP",
          boardCards: ["♥Q", "**", "**"],
          seats: [
            {
              index: 5,
              name: "Player6",
              status: "PLAYING",
              bankroll: 99,
              inPotAmount: 1,
              isTurn: true,
              cards: ["♣A", "♣K"],
            },
          ],
        }}
        onAction={async () => {}}
        onStartHand={async () => {}}
        onTakeSeat={async () => {}}
      />,
    );

    expect(screen.queryByText(/leave this table/i)).not.toBeInTheDocument();

    fireEvent.click(screen.getByRole("button", { name: /back to lobby/i }));

    expect(screen.getByText(/leave this table/i)).toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: /stay here/i }),
    ).toBeInTheDocument();

    const leaveLink = screen.getByRole("link", { name: /leave room/i });
    expect(leaveLink).toHaveAttribute("href", "/");

    fireEvent.click(screen.getByRole("button", { name: /stay here/i }));

    expect(screen.queryByText(/leave this table/i)).not.toBeInTheDocument();
  });

  it("shows current-round seat actions derived from structured room events", () => {
    render(
      <RoomPage
        snapshot={{
          roomId: "room-001",
          roomName: "Table 1",
          status: "awaiting_action",
          viewerRole: "player",
          handNumber: 9,
          humanSeat: 5,
          smallBlind: 1,
          pot: 6,
          currentAmount: 2,
          round: "FLOP",
          boardCards: ["♥Q", "♣J", "**"],
          seats: [
            {
              index: 0,
              name: "Player1",
              status: "PLAYING",
              bankroll: 98,
              inPotAmount: 2,
              isTurn: false,
              cards: ["**", "**"],
            },
            {
              index: 3,
              name: "Player4",
              status: "OUT",
              bankroll: 100,
              inPotAmount: 0,
              isTurn: false,
              cards: ["**", "**"],
            },
            {
              index: 5,
              name: "Player6",
              status: "PLAYING",
              bankroll: 99,
              inPotAmount: 1,
              isTurn: true,
              cards: ["♣A", "♣K"],
            },
          ],
          pendingAction: {
            token: "turn-1",
            seatIndex: 5,
            minAmount: 1,
            maxAmount: 99,
            canCheck: false,
            canCall: true,
            canBet: true,
            canFold: true,
            canAllIn: true,
          },
          events: [
            {
              kind: "round_start",
              message: "preflop opened",
              round: "PREFLOP",
              handNumber: 9,
            },
            {
              kind: "player_action",
              message: "seat 1 called 2",
              round: "PREFLOP",
              handNumber: 9,
              seatIndex: 0,
              actionType: "CALL",
              amount: 2,
            },
            {
              kind: "round_start",
              message: "flop opened",
              round: "FLOP",
              handNumber: 9,
            },
            {
              kind: "player_action",
              message: "seat 1 bet 4",
              round: "FLOP",
              handNumber: 9,
              seatIndex: 0,
              actionType: "BET",
              amount: 4,
            },
            {
              kind: "player_action",
              message: "seat 4 folded",
              round: "FLOP",
              handNumber: 9,
              seatIndex: 3,
              actionType: "FOLD",
              amount: 0,
            },
          ],
        }}
        onAction={async () => {}}
        onStartHand={async () => {}}
        onTakeSeat={async () => {}}
      />,
    );

    expect(screen.getByText("Bet 4")).toBeInTheDocument();
    expect(screen.getByText("Folded")).toBeInTheDocument();
    expect(screen.getByText("Latest event")).toBeInTheDocument();
    expect(screen.getAllByText("Player4 folds").length).toBeGreaterThanOrEqual(
      2,
    );
  });

  it("surfaces blind, deal, and pot collection events as table action cues", () => {
    const { container, rerender } = render(
      <RoomPage
        snapshot={{
          roomId: "room-001",
          roomName: "Table 1",
          status: "running",
          viewerRole: "spectator",
          humanSeat: 5,
          handNumber: 3,
          smallBlind: 1,
          pot: 3,
          currentAmount: 2,
          round: "PREFLOP",
          boardCards: [],
          seats: [
            {
              index: 0,
              name: "Player1",
              status: "PLAYING",
              bankroll: 99,
              inPotAmount: 1,
              isTurn: false,
              cards: ["**", "**"],
            },
          ],
          events: [
            {
              kind: "hole_cards_dealt",
              message: "hole cards dealt",
              handNumber: 3,
              round: "PREFLOP",
            },
            {
              kind: "blind_posted",
              message: "small blind posted",
              handNumber: 3,
              round: "PREFLOP",
              seatIndex: 0,
              actionType: "SMALL_BLIND",
              amount: 1,
            },
          ],
        }}
        onAction={async () => {}}
        onStartHand={async () => {}}
        onTakeSeat={async () => {}}
      />,
    );

    expect(screen.getByText("Latest event")).toBeInTheDocument();
    expect(
      screen.getAllByText("Player1 posts small blind 1").length,
    ).toBeGreaterThanOrEqual(2);
    expect(
      container.querySelector(".table-live-layout--blind-posted"),
    ).toBeInTheDocument();
    expect(screen.queryByText("Player1 to Pot +1")).not.toBeInTheDocument();
    expect(container.querySelector(".chip-flow-cue")).not.toBeInTheDocument();

    rerender(
      <RoomPage
        snapshot={{
          roomId: "room-001",
          roomName: "Table 1",
          status: "hand_finished",
          viewerRole: "spectator",
          humanSeat: 5,
          handNumber: 3,
          smallBlind: 1,
          pot: 0,
          currentAmount: 0,
          round: "FINISH",
          boardCards: ["♥Q", "♣J", "♦10", "♠2", "♣3"],
          seats: [
            {
              index: 0,
              name: "Player1",
              status: "PLAYING",
              bankroll: 103,
              inPotAmount: 0,
              isTurn: false,
              isWinner: true,
              netChange: 3,
              cards: ["♥A", "♠K"],
            },
          ],
          events: [
            {
              kind: "pot_collected",
              message: "pot collected",
              handNumber: 3,
              round: "FINISH",
              amount: 3,
            },
            {
              kind: "hand_finish",
              message: "hand 3 finished",
              handNumber: 3,
            },
          ],
        }}
        onAction={async () => {}}
        onStartHand={async () => {}}
        onTakeSeat={async () => {}}
      />,
    );

    expect(screen.getAllByText("Player1 wins 3").length).toBeGreaterThanOrEqual(
      2,
    );
    expect(
      container.querySelector(".table-live-layout--pot-collected"),
    ).toBeInTheDocument();
    expect(screen.queryByText("Pot to Player1 +3")).not.toBeInTheDocument();
    expect(container.querySelector(".chip-flow-cue")).not.toBeInTheDocument();
  });

  it("shows dealt streets in the latest table event cue", () => {
    render(
      <RoomPage
        snapshot={{
          roomId: "room-001",
          roomName: "Table 1",
          status: "running",
          viewerRole: "spectator",
          humanSeat: 5,
          handNumber: 4,
          smallBlind: 1,
          pot: 12,
          currentAmount: 0,
          round: "FLOP",
          boardCards: ["♥Q", "♣J", "♦10", "**", "**"],
          seats: [
            {
              index: 0,
              name: "Player1",
              status: "PLAYING",
              bankroll: 98,
              inPotAmount: 2,
              isTurn: false,
              cards: ["**", "**"],
            },
          ],
          events: [
            {
              kind: "round_start",
              message: "flop opened",
              handNumber: 4,
              round: "FLOP",
            },
          ],
        }}
        onAction={async () => {}}
        onStartHand={async () => {}}
        onTakeSeat={async () => {}}
      />,
    );

    expect(screen.getByText("Latest event")).toBeInTheDocument();
    expect(screen.getAllByText("Flop dealt").length).toBeGreaterThanOrEqual(2);
  });

  it("keeps spectators out of the live action controls even when a hand is awaiting input", () => {
    render(
      <RoomPage
        snapshot={{
          roomId: "room-001",
          roomName: "Table 1",
          status: "awaiting_action",
          viewerRole: "spectator",
          humanSeat: 5,
          handNumber: 3,
          smallBlind: 1,
          pot: 6,
          currentAmount: 2,
          round: "FLOP",
          boardCards: ["♥Q", "**", "**"],
          seats: [
            {
              index: 0,
              name: "Player1",
              status: "PLAYING",
              bankroll: 98,
              inPotAmount: 2,
              isTurn: false,
              cards: ["**", "**"],
            },
            {
              index: 5,
              name: "Player6",
              status: "PLAYING",
              bankroll: 99,
              inPotAmount: 1,
              isTurn: true,
              cards: ["**", "**"],
            },
          ],
          pendingAction: {
            token: "turn-1",
            seatIndex: 5,
            minAmount: 1,
            maxAmount: 99,
            canCheck: false,
            canCall: true,
            canBet: true,
            canFold: true,
            canAllIn: true,
          },
        }}
        onAction={async () => {}}
        onStartHand={async () => {}}
        onTakeSeat={async () => {}}
      />,
    );

    expect(
      screen.queryByRole("button", { name: /confirm bet/i }),
    ).not.toBeInTheDocument();
    expect(
      screen.queryByRole("button", { name: /spectate/i }),
    ).not.toBeInTheDocument();
    expect(
      screen.getByText(
        /your cards are hidden because this room is currently in spectator mode/i,
      ),
    ).toBeInTheDocument();
  });

  it("disables the start-hand control while a hand is already in progress", () => {
    render(
      <RoomPage
        snapshot={{
          roomId: "room-001",
          roomName: "Table 1",
          status: "awaiting_action",
          viewerRole: "spectator",
          humanSeat: 5,
          handNumber: 4,
          smallBlind: 1,
          pot: 3,
          currentAmount: 2,
          round: "PREFLOP",
          boardCards: ["**", "**", "**", "**", "**"],
          seats: [
            {
              index: 5,
              name: "Player6",
              status: "PLAYING",
              bankroll: 98,
              inPotAmount: 2,
              isTurn: true,
              cards: ["**", "**"],
            },
          ],
          pendingAction: {
            token: "turn-1",
            seatIndex: 5,
            minAmount: 2,
            maxAmount: 98,
            canCheck: false,
            canCall: true,
            canBet: true,
            canFold: true,
            canAllIn: true,
          },
        }}
        onAction={async () => {}}
        onStartHand={async () => {}}
        onTakeSeat={async () => {}}
      />,
    );

    const startHandButton = screen.getByRole("button", {
      name: /hand in progress/i,
    });
    expect(startHandButton).toBeDisabled();
    expect(startHandButton).not.toHaveClass("primary");
    expect(startHandButton).toHaveClass("secondary");
  });

  it("leaves the room, releases the seat, and clears the stored viewer session", async () => {
    window.localStorage.setItem("poker.viewerSeat.room-001", "5");
    window.localStorage.setItem("poker.viewerToken.room-001", "viewer-token-1");

    const playerSnapshot: RoomSnapshot = {
      roomId: "room-001",
      roomName: "Table 1",
      status: "awaiting_action",
      viewerRole: "player",
      humanSeat: 5,
      handNumber: 3,
      smallBlind: 1,
      pot: 6,
      currentAmount: 2,
      round: "FLOP",
      boardCards: ["♥Q", "**", "**"],
      seats: [
        {
          index: 5,
          name: "Player6",
          status: "PLAYING",
          bankroll: 99,
          inPotAmount: 1,
          isTurn: true,
          cards: ["♣A", "♣K"],
        },
      ],
      pendingAction: {
        token: "turn-1",
        seatIndex: 5,
        minAmount: 1,
        maxAmount: 99,
        canCheck: false,
        canCall: true,
        canBet: true,
        canFold: true,
        canAllIn: true,
      },
    };
    const spectatorSnapshot: RoomSnapshot = {
      roomId: "room-001",
      roomName: "Table 1",
      status: "running",
      viewerRole: "spectator",
      humanSeat: 5,
      handNumber: 3,
      smallBlind: 1,
      pot: 6,
      currentAmount: 2,
      round: "FLOP",
      boardCards: ["♥Q", "**", "**"],
      seats: [
        {
          index: 5,
          name: "Player6",
          status: "PLAYING",
          bankroll: 99,
          inPotAmount: 1,
          isTurn: true,
          cards: ["**", "**"],
        },
      ],
    };

    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce({
        ok: true,
        json: async () => playerSnapshot,
      })
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          roomId: "room-001",
        }),
      });

    vi.stubGlobal("fetch", fetchMock);
    vi.stubGlobal("WebSocket", MockWebSocket as unknown as typeof WebSocket);

    render(
      <MemoryRouter initialEntries={["/rooms/room-001"]}>
        <Routes>
          <Route path="/" element={<div>Lobby landing</div>} />
          <Route path="/rooms/:roomId" element={<RoomRoute />} />
        </Routes>
      </MemoryRouter>,
    );

    await waitFor(() =>
      expect(screen.getByText("Table 1")).toBeInTheDocument(),
    );

    expect(fetchMock).toHaveBeenNthCalledWith(
      1,
      "/api/rooms/room-001?viewerSeat=5&viewerToken=viewer-token-1",
    );
    expect(MockWebSocket.instances[0]?.url).toContain(
      "/ws/rooms/room-001?viewerSeat=5&viewerToken=viewer-token-1",
    );

    fireEvent.click(screen.getByRole("button", { name: /back to lobby/i }));
    fireEvent.click(screen.getByRole("link", { name: /leave room/i }));

    await waitFor(() =>
      expect(fetchMock).toHaveBeenNthCalledWith(
        2,
        "/api/rooms/room-001/leave",
        expect.objectContaining({
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ viewerToken: "viewer-token-1" }),
        }),
      ),
    );
    await waitFor(() =>
      expect(screen.getByText("Lobby landing")).toBeInTheDocument(),
    );

    expect(window.localStorage.getItem("poker.viewerSeat.room-001")).toBeNull();
    expect(
      window.localStorage.getItem("poker.viewerToken.room-001"),
    ).toBeNull();
    expect(MockWebSocket.instances[0]?.close).toHaveBeenCalled();
  });

  it("plays queued websocket snapshots with a perceptible action delay", async () => {
    window.localStorage.setItem("poker.viewerSeat.room-001", "5");
    window.localStorage.setItem("poker.viewerToken.room-001", "viewer-token-1");

    const initialSnapshot: RoomSnapshot = {
      roomId: "room-001",
      roomName: "Table 1",
      status: "running",
      viewerRole: "player",
      humanSeat: 5,
      playerCount: 6,
      handNumber: 4,
      smallBlind: 1,
      pot: 3,
      currentAmount: 2,
      round: "PREFLOP",
      boardCards: [],
      seats: [
        {
          index: 0,
          name: "Player1",
          status: "PLAYING",
          bankroll: 98,
          inPotAmount: 2,
          isTurn: false,
          cards: ["**", "**"],
        },
        {
          index: 5,
          name: "Player6",
          status: "PLAYING",
          bankroll: 99,
          inPotAmount: 1,
          isTurn: true,
          cards: ["♣A", "♣K"],
        },
      ],
      pendingAction: {
        token: "turn-1",
        seatIndex: 5,
        minAmount: 1,
        maxAmount: 99,
        canCheck: false,
        canCall: true,
        canBet: true,
        canFold: true,
        canAllIn: true,
      },
      events: [],
      version: 1,
    };

    const firstActionSnapshot: RoomSnapshot = {
      ...initialSnapshot,
      pot: 5,
      seats: [
        {
          index: 0,
          name: "Player1",
          status: "PLAYING",
          bankroll: 96,
          inPotAmount: 4,
          isTurn: false,
          cards: ["**", "**"],
        },
        {
          index: 5,
          name: "Player6",
          status: "PLAYING",
          bankroll: 99,
          inPotAmount: 1,
          isTurn: true,
          cards: ["♣A", "♣K"],
        },
      ],
      events: [
        {
          kind: "player_action",
          message: "seat 0 called 2",
          handNumber: 4,
          round: "PREFLOP",
          seatIndex: 0,
          actionType: "CALL",
          amount: 2,
        },
      ],
      version: 2,
    };

    const secondActionSnapshot: RoomSnapshot = {
      ...firstActionSnapshot,
      status: "hand_finished",
      pendingAction: undefined,
      seats: [
        {
          index: 0,
          name: "Player1",
          status: "PLAYING",
          bankroll: 96,
          inPotAmount: 4,
          isTurn: false,
          cards: ["**", "**"],
        },
        {
          index: 5,
          name: "Player6",
          status: "PLAYING",
          bankroll: 104,
          inPotAmount: 0,
          isTurn: false,
          isWinner: true,
          netChange: 5,
          bestHand: "One pair",
          cards: ["♣A", "♣K"],
        },
      ],
      events: [
        ...firstActionSnapshot.events!,
        {
          kind: "player_action",
          message: "seat 5 called 2",
          handNumber: 4,
          round: "PREFLOP",
          seatIndex: 5,
          actionType: "CALL",
          amount: 2,
        },
      ],
      version: 3,
    };

    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => initialSnapshot,
    });

    vi.stubGlobal("fetch", fetchMock);
    vi.stubGlobal("WebSocket", MockWebSocket as unknown as typeof WebSocket);

    render(
      <MemoryRouter initialEntries={["/rooms/room-001"]}>
        <Routes>
          <Route path="/" element={<div>Lobby landing</div>} />
          <Route path="/rooms/:roomId" element={<RoomRoute />} />
        </Routes>
      </MemoryRouter>,
    );

    await waitFor(() =>
      expect(screen.getByText("Table 1")).toBeInTheDocument(),
    );

    vi.useFakeTimers();
    try {
      await act(async () => {
        MockWebSocket.instances[0]?.onmessage?.({
          data: JSON.stringify(firstActionSnapshot),
        } as MessageEvent);
        MockWebSocket.instances[0]?.onmessage?.({
          data: JSON.stringify(secondActionSnapshot),
        } as MessageEvent);
      });

      expect(screen.queryByText("Player6 calls 2")).not.toBeInTheDocument();
      expect(screen.queryByText("Player6 wins +5")).not.toBeInTheDocument();

      await act(async () => {
        await vi.advanceTimersByTimeAsync(ACTION_PLAYBACK_DELAY_MS);
      });

      expect(screen.getAllByText("Player1 calls 2").length).toBeGreaterThanOrEqual(
        2,
      );
      expect(screen.queryByText("Player6 calls 2")).not.toBeInTheDocument();

      await act(async () => {
        await vi.advanceTimersByTimeAsync(ACTION_PLAYBACK_DELAY_MS);
      });

      expect(screen.getAllByText("Player6 calls 2").length).toBeGreaterThanOrEqual(
        2,
      );
      expect(screen.getByText("Player6 wins +5")).toBeInTheDocument();
    } finally {
      vi.useRealTimers();
    }
  });

  it("shows connection state and resubscribes after the room socket closes", async () => {
    window.localStorage.setItem("poker.viewerSeat.room-001", "5");
    window.localStorage.setItem("poker.viewerToken.room-001", "viewer-token-1");

    const initialSnapshot: RoomSnapshot = {
      roomId: "room-001",
      roomName: "Table 1",
      status: "running",
      viewerRole: "player",
      humanSeat: 5,
      playerCount: 6,
      handNumber: 4,
      smallBlind: 1,
      pot: 3,
      currentAmount: 2,
      round: "PREFLOP",
      boardCards: [],
      seats: [
        {
          index: 5,
          name: "Player6",
          status: "PLAYING",
          bankroll: 99,
          inPotAmount: 1,
          isTurn: true,
          cards: ["♣A", "♣K"],
        },
      ],
      version: 1,
    };
    const resumedSnapshot = {
      ...initialSnapshot,
      pot: 5,
      version: 2,
    };

    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce({
        ok: true,
        json: async () => initialSnapshot,
      })
      .mockResolvedValueOnce({
        ok: true,
        json: async () => resumedSnapshot,
      });

    vi.stubGlobal("fetch", fetchMock);
    vi.stubGlobal("WebSocket", MockWebSocket as unknown as typeof WebSocket);

    render(
      <MemoryRouter initialEntries={["/rooms/room-001"]}>
        <Routes>
          <Route path="/" element={<div>Lobby landing</div>} />
          <Route path="/rooms/:roomId" element={<RoomRoute />} />
        </Routes>
      </MemoryRouter>,
    );

    await waitFor(() =>
      expect(screen.getByText("Table 1")).toBeInTheDocument(),
    );

    act(() => {
      MockWebSocket.instances[0]?.onopen?.();
    });

    expect(screen.getByText("live connection")).toBeInTheDocument();

    vi.useFakeTimers();
    try {
      act(() => {
        MockWebSocket.instances[0]?.onclose?.();
      });

      expect(screen.getByText("reconnecting")).toBeInTheDocument();

      await act(async () => {
        await vi.advanceTimersByTimeAsync(500);
      });
      await act(async () => {
        await Promise.resolve();
        await Promise.resolve();
      });

      expect(MockWebSocket.instances).toHaveLength(2);
      expect(fetchMock).toHaveBeenNthCalledWith(
        2,
        "/api/rooms/room-001?viewerSeat=5&viewerToken=viewer-token-1",
      );
    } finally {
      vi.useRealTimers();
    }
  });
});
