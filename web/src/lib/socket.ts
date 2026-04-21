import type { RoomSnapshot } from "./types";

function buildSocketURL(
  roomId: string,
  viewerSeat?: number,
  viewerToken?: string,
): string {
  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  const url = new URL(
    `${protocol}//${window.location.host}/ws/rooms/${roomId}`,
  );
  if (viewerSeat !== undefined) {
    url.searchParams.set("viewerSeat", String(viewerSeat));
    if (viewerToken) {
      url.searchParams.set("viewerToken", viewerToken);
    }
  }
  return url.toString();
}

export function subscribeRoom(
  roomId: string,
  viewerSeat: number | undefined,
  viewerToken: string | undefined,
  onSnapshot: (snapshot: RoomSnapshot) => void,
  onError?: (message: string) => void,
): () => void {
  const socket = new WebSocket(buildSocketURL(roomId, viewerSeat, viewerToken));

  socket.onmessage = (event) => {
    try {
      const payload = JSON.parse(event.data) as
        | RoomSnapshot
        | { error?: string };
      if ("error" in payload && payload.error) {
        onError?.(payload.error);
        return;
      }
      onSnapshot(payload as RoomSnapshot);
    } catch (err) {
      onError?.(
        err instanceof Error ? err.message : "failed to parse room update",
      );
    }
  };

  socket.onerror = () => {
    onError?.("room socket disconnected");
  };

  return () => {
    socket.close();
  };
}
