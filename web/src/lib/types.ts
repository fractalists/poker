export type PendingAction = {
  token: string;
  seatIndex: number;
  minAmount: number;
  minBetAmount?: number;
  maxAmount: number;
  canCheck: boolean;
  canCall: boolean;
  canBet: boolean;
  canFold: boolean;
  canAllIn: boolean;
};

export type SeatSnapshot = {
  index: number;
  name: string;
  status: string;
  bankroll: number;
  inPotAmount: number;
  isTurn: boolean;
  isWinner?: boolean;
  netChange?: number;
  bestHand?: string;
  cards: string[];
};

export type RoomSnapshot = {
  roomId: string;
  roomName: string;
  humanSeat?: number;
  playerCount?: number;
  status: string;
  viewerRole?: "player" | "spectator";
  handNumber: number;
  smallBlind: number;
  pot?: number;
  currentAmount?: number;
  round?: string;
  boardCards?: string[];
  seats: SeatSnapshot[];
  pendingAction?: PendingAction;
  events?: RoomEvent[];
  version?: number;
};

export type RoomEvent = {
  kind: string;
  message: string;
  handNumber?: number;
  round?: string;
  seatIndex?: number;
  actionType?: string;
  amount?: number;
};

export type ViewerSession = {
  roomId: string;
  viewerSeat?: number;
  viewerToken?: string;
};

export type ActionSubmission = {
  token: string;
  actionType: "BET" | "CALL" | "FOLD" | "ALL_IN";
  amount: number;
};
