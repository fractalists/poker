import type { ActionSubmission, RoomSnapshot, ViewerSession } from "./types";

function buildRoomURL(
  roomId: string,
  viewerSeat?: number,
  viewerToken?: string,
): string {
  const url = new URL(`/api/rooms/${roomId}`, window.location.origin);
  if (viewerSeat !== undefined) {
    url.searchParams.set("viewerSeat", String(viewerSeat));
    if (viewerToken) {
      url.searchParams.set("viewerToken", viewerToken);
    }
  }
  return `${url.pathname}${url.search}`;
}

async function parseJSON<T>(response: Response): Promise<T> {
  if (!response.ok) {
    throw new Error(
      (await response.text()) || `request failed: ${response.status}`,
    );
  }
  return response.json() as Promise<T>;
}

async function ensureOK(response: Response): Promise<void> {
  if (!response.ok) {
    throw new Error(
      (await response.text()) || `request failed: ${response.status}`,
    );
  }
}

export async function listRooms(): Promise<RoomSnapshot[]> {
  return parseJSON<RoomSnapshot[]>(await fetch("/api/rooms"));
}

export async function createRoom(input: {
  name: string;
  smallBlind: number;
  startingBankroll: number;
  humanSeat: number;
  playerCount: number;
  aiStyle: string;
}): Promise<RoomSnapshot> {
  return parseJSON<RoomSnapshot>(
    await fetch("/api/rooms", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(input),
    }),
  );
}

export async function getRoom(
  roomId: string,
  viewerSeat?: number,
  viewerToken?: string,
): Promise<RoomSnapshot> {
  return parseJSON<RoomSnapshot>(
    await fetch(buildRoomURL(roomId, viewerSeat, viewerToken)),
  );
}

export async function takeSeat(
  roomId: string,
  seatIndex: number,
): Promise<ViewerSession> {
  return parseJSON<ViewerSession>(
    await fetch(`/api/rooms/${roomId}/seat`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ seatIndex }),
    }),
  );
}

export async function leaveRoom(
  roomId: string,
  viewerToken?: string,
): Promise<void> {
  const init: RequestInit = { method: "POST" };
  if (viewerToken) {
    init.headers = {
      "Content-Type": "application/json",
    };
    init.body = JSON.stringify({ viewerToken });
  }

  await ensureOK(await fetch(`/api/rooms/${roomId}/leave`, init));
}

export async function startHand(roomId: string): Promise<void> {
  await ensureOK(
    await fetch(`/api/rooms/${roomId}/start`, {
      method: "POST",
    }),
  );
}

export async function submitAction(
  roomId: string,
  input: ActionSubmission,
  viewerToken?: string,
): Promise<void> {
  await ensureOK(
    await fetch(`/api/rooms/${roomId}/actions`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ ...input, viewerToken }),
    }),
  );
}
