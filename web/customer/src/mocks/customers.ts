import type {
  ActivitiesQuery,
  ActivitiesResponse,
  Bandwidth,
  ChangePasswordInput,
  CiscoSetup,
  DailyTraffic,
  LoginData,
  LoginResponse,
  OnlineUserSession,
  SummaryResponse,
} from "@/api/generated";
import { ApiError } from "@/api/http";

const scenario = () => import.meta.env.VITE_CUSTOMER_MOCK_SCENARIO || "success";
async function ready(): Promise<void> {
  await new Promise((resolve) => setTimeout(resolve, 150));
  if (scenario() === "error") throw new ApiError("Service unavailable", 503);
}
const summary: SummaryResponse = {
  ocserv_user: {
    username: "customer",
    owner: "Customer",
    expiry_mode: "fixed",
    expire_at: "2027-01-01",
    traffic_type: "MonthlyRxTx",
    traffic_size: 100,
    running_rx: 12,
    running_tx: 4,
    certificate_available: true,
    certificate_enabled: true,
    is_locked: false,
  },
  usage: {
    date_start: "2026-09-01",
    date_end: "2026-09-30",
    bandwidths: { rx: 12, tx: 4 },
  },
};
const sessions: OnlineUserSession[] = [
  {
    ID: 1,
    Device: "vpn0",
    IPv4: "10.10.0.2",
    "Session started at": "2026-09-11T10:00:00Z",
    vhost: "default",
    "Average RX": "1.2 MiB/s",
    "Average TX": "420 KiB/s",
  },
];
export const mockCustomers = {
  async login(request: LoginData): Promise<LoginResponse> {
    await ready();
    if (request.username === "invalid")
      throw new ApiError("Invalid credentials", 400);
    return {
      token: "mock-customer-token",
      expires_at: new Date(Date.now() + 7 * 60 * 60 * 1000).toISOString(),
      user: { ...summary.ocserv_user, username: request.username },
    };
  },
  async summary(): Promise<SummaryResponse> {
    await ready();
    return structuredClone(summary);
  },
  async sessions(): Promise<OnlineUserSession[]> {
    await ready();
    return scenario() === "empty" ? [] : structuredClone(sessions);
  },
  async activities(params: ActivitiesQuery): Promise<ActivitiesResponse> {
    await ready();
    const result =
      scenario() === "empty"
        ? []
        : [
            {
              username: "customer",
              event: "handshake" as const,
              message: "Session established",
              ip: "192.0.2.10",
              created_at: "2026-09-11T10:00:00Z",
            },
          ];
    return {
      meta: {
        page: params.page ?? 1,
        size: params.size ?? 20,
        total_records: result.length,
      },
      result,
    };
  },
  async stats(): Promise<DailyTraffic[]> {
    await ready();
    return scenario() === "empty"
      ? []
      : [{ date: "2026-09-11", rx: 1.2, tx: 0.4 }];
  },
  async bandwidth(): Promise<Bandwidth> {
    await ready();
    return { rx: 12, tx: 4 };
  },
  async ciscoSetup(): Promise<CiscoSetup> {
    await ready();
    return {
      connection_name: "Ocserv VPN",
      server_address: "vpn.example.test",
      server_port: 443,
      certificate_password: "customer",
      certificate_import_uri: "cisco-secure-client://import",
      connection_create_uri: "cisco-secure-client://connect",
    };
  },
  async certificate(): Promise<Blob> {
    await ready();
    return new Blob(["mock certificate"], { type: "application/x-pkcs12" });
  },
  async accept(): Promise<void> {
    await ready();
  },
  async changePassword(request: ChangePasswordInput): Promise<void> {
    await ready();
    if (request.password.length < 2)
      throw new ApiError("Password must contain at least 2 characters", 400);
  },
};
