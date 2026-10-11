import { afterEach, describe, expect, it, vi } from "vitest";

vi.mock("service/http", () => ({ fetch: vi.fn() }));

import { fetch as apiFetch } from "service/http";

import { useServicesStore } from "./ServicesContext";

type PendingRequest = { resolve: (value: unknown) => void };

describe("fetchServiceDetail", () => {
	afterEach(() => {
		vi.clearAllMocks();
		useServicesStore.setState({ isLoading: false, serviceDetail: null });
	});

	it("keeps the newest service detail when requests finish out of order", async () => {
		const pending: PendingRequest[] = [];
		vi.mocked(apiFetch).mockImplementation(
			() =>
				new Promise((resolve) => {
					pending.push({ resolve });
				}),
		);

		const first = useServicesStore.getState().fetchServiceDetail(1);
		const second = useServicesStore.getState().fetchServiceDetail(2);

		pending[1].resolve({ id: 2, name: "new" });
		expect(await second).toMatchObject({ id: 2 });
		pending[0].resolve({ id: 1, name: "old" });
		expect(await first).toBeNull();

		expect(useServicesStore.getState().serviceDetail).toMatchObject({ id: 2 });
		expect(useServicesStore.getState().isLoading).toBe(false);
	});
});

describe("fetchServices", () => {
	afterEach(() => {
		vi.clearAllMocks();
		useServicesStore.setState({ services: [], total: 0, isLoading: false });
	});

	it("requests every service unless a caller asks for a page", async () => {
		vi.mocked(apiFetch).mockResolvedValue({ services: [], total: 0 });

		await useServicesStore.getState().fetchServices();
		expect(vi.mocked(apiFetch).mock.calls.at(-1)?.[1]).toMatchObject({
			query: { limit: "none", offset: 0 },
		});

		await useServicesStore.getState().fetchServices({ limit: 5, offset: 10 });
		expect(vi.mocked(apiFetch).mock.calls.at(-1)?.[1]).toMatchObject({
			query: { limit: 5, offset: 10 },
		});
	});
});
