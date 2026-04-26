import { act, fireEvent, render, screen, waitFor } from "@testing-library/react";
import { vi } from "vitest";

import { LobbyPage } from "./LobbyPage";

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

describe("LobbyPage", () => {
  beforeEach(() => {
    vi.stubGlobal("WebSocket", MockWebSocket as unknown as typeof WebSocket);
  });

  afterEach(() => {
    MockWebSocket.instances = [];
    vi.restoreAllMocks();
    vi.unstubAllGlobals();
    window.localStorage.clear();
  });

  it("renders rooms from the service and exposes create-room controls", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: true,
        json: async () => [
          {
            roomId: "room-001",
            roomName: "Table 1",
            status: "waiting",
            handNumber: 0,
            smallBlind: 1,
            playerCount: 6,
            seats: [],
          },
        ],
      }),
    );

    const { container } = render(<LobbyPage />);

    await waitFor(() =>
      expect(screen.getByText("Table 1")).toBeInTheDocument(),
    );
    expect(
      screen.getByRole("button", { name: /create room/i }),
    ).toBeInTheDocument();
    expect(screen.getByLabelText(/total players/i)).toBeInTheDocument();
    expect(screen.getByRole("option", { name: /10 players/i })).toBeInTheDocument();
    expect(screen.getByRole("option", { name: /mixed ai/i })).toBeInTheDocument();
    expect(screen.queryByRole("option", { name: /random seats/i })).not.toBeInTheDocument();
    expect(screen.getByRole("option", { name: /gto/i })).toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: /spectate room/i }),
    ).toBeInTheDocument();
    const workGrid = container.querySelector(".work-grid");
    expect(workGrid?.firstElementChild).toHaveClass("room-strip");
    expect(workGrid?.lastElementChild).toHaveClass("create-panel");
  });

  it("creates a room with the selected player count and auto claims a random human seat", async () => {
    const navigateToRoom = vi.fn();
    vi.spyOn(Math, "random").mockReturnValue(0.25);
    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce({
        ok: true,
        json: async () => [],
      })
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          roomId: "room-002",
          roomName: "Table 2",
          humanSeat: 1,
          status: "waiting",
          handNumber: 0,
          smallBlind: 1,
          playerCount: 4,
          aiStyle: "aggressive",
          seats: [],
        }),
      })
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          roomId: "room-002",
          viewerSeat: 1,
          viewerToken: "viewer-token-1",
        }),
      });

    vi.stubGlobal("fetch", fetchMock);

    render(<LobbyPage navigateToRoom={navigateToRoom} />);

    await waitFor(() =>
      expect(
        screen.getByRole("button", { name: /create room/i }),
      ).toBeInTheDocument(),
    );
    fireEvent.change(screen.getByLabelText(/total players/i), {
      target: { value: "4" },
    });
    fireEvent.change(screen.getByLabelText(/ai style/i), {
      target: { value: "aggressive" },
    });
    fireEvent.click(screen.getByRole("button", { name: /create room/i }));

    await waitFor(() =>
      expect(navigateToRoom).toHaveBeenCalledWith("room-002"),
    );
    expect(window.localStorage.getItem("poker.viewerSeat.room-002")).toBe("1");
    expect(window.localStorage.getItem("poker.viewerToken.room-002")).toBe(
      "viewer-token-1",
    );
    expect(fetchMock).toHaveBeenNthCalledWith(
      2,
      "/api/rooms",
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify({
          name: "Table 1",
          smallBlind: 1,
          startingBankroll: 100,
          humanSeat: 1,
          playerCount: 4,
          aiStyle: "aggressive",
        }),
      }),
    );
    expect(fetchMock).toHaveBeenNthCalledWith(
      3,
      "/api/rooms/room-002/seat",
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify({ seatIndex: 1 }),
      }),
    );
  });

  it("can create a room with mixed per-seat AI selection", async () => {
    const navigateToRoom = vi.fn();
    vi.spyOn(Math, "random").mockReturnValue(0.25);
    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce({
        ok: true,
        json: async () => [],
      })
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          roomId: "room-003",
          roomName: "Table 1",
          humanSeat: 1,
          status: "waiting",
          handNumber: 0,
          smallBlind: 1,
          playerCount: 4,
          aiStyle: "random",
          seats: [],
        }),
      })
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          roomId: "room-003",
          viewerSeat: 1,
          viewerToken: "viewer-token-1",
        }),
      });

    vi.stubGlobal("fetch", fetchMock);

    render(<LobbyPage navigateToRoom={navigateToRoom} />);

    await waitFor(() =>
      expect(
        screen.getByRole("button", { name: /create room/i }),
      ).toBeInTheDocument(),
    );
    fireEvent.change(screen.getByLabelText(/total players/i), {
      target: { value: "4" },
    });
    fireEvent.click(screen.getByRole("button", { name: /create room/i }));

    await waitFor(() =>
      expect(navigateToRoom).toHaveBeenCalledWith("room-003"),
    );
    expect(fetchMock).toHaveBeenNthCalledWith(
      2,
      "/api/rooms",
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify({
          name: "Table 1",
          smallBlind: 1,
          startingBankroll: 100,
          humanSeat: 1,
          playerCount: 4,
          aiStyle: "random",
        }),
      }),
    );
  });

  it("keeps the lobby room list in sync from the rooms websocket", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: true,
        json: async () => [],
      }),
    );

    render(<LobbyPage />);

    await waitFor(() =>
      expect(screen.getByText(/no rooms yet/i)).toBeInTheDocument(),
    );
    expect(MockWebSocket.instances[0]?.url).toContain("/ws/rooms");

    act(() => {
      MockWebSocket.instances[0]?.onmessage?.({
        data: JSON.stringify([
          {
            roomId: "room-099",
            roomName: "Socket Table",
            status: "waiting",
            handNumber: 0,
            smallBlind: 1,
            playerCount: 6,
            aiStyle: "mixed",
            seats: [],
          },
        ]),
      } as MessageEvent);
    });

    expect(screen.getByText("Socket Table")).toBeInTheDocument();
    expect(screen.getByText("AI mixed")).toBeInTheDocument();
  });

  it("enters a room in pure spectator mode from the lobby and clears any saved viewer session", async () => {
    const navigateToRoom = vi.fn();
    window.localStorage.setItem("poker.viewerSeat.room-001", "5");
    window.localStorage.setItem("poker.viewerToken.room-001", "viewer-token-1");
    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce({
        ok: true,
        json: async () => [
          {
            roomId: "room-001",
            roomName: "Table 1",
            status: "waiting",
            handNumber: 0,
            smallBlind: 1,
            playerCount: 6,
            seats: [],
          },
        ],
      })
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          roomId: "room-001",
        }),
      });

    vi.stubGlobal(
      "fetch",
      fetchMock,
    );

    render(<LobbyPage navigateToRoom={navigateToRoom} />);

    await waitFor(() =>
      expect(screen.getByText("Table 1")).toBeInTheDocument(),
    );

    fireEvent.click(screen.getByRole("button", { name: /spectate room/i }));

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
      expect(navigateToRoom).toHaveBeenCalledWith("room-001"),
    );
    expect(window.localStorage.getItem("poker.viewerSeat.room-001")).toBeNull();
    expect(window.localStorage.getItem("poker.viewerToken.room-001")).toBeNull();
  });
});
