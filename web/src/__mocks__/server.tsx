import { rest } from "msw";
import { setupServer } from "msw/node";

const server = setupServer(
	rest.get("/status", (req, res, ctx) => {
		return res(
			ctx.json({
				status: "ok",
				timestamp: "2021-08-16T18:16:29.113226259Z",
				version: "1.0.1"
			})
		);
	}),

	rest.get("/secrets", (req, res, ctx) => {
		return res(
			ctx.json({
				secret: "the eagle flies at midnight",
				is_base64: false,
				created: "2021-08-16T16:48:14.584642555Z",
				accesses: 1,
				destroyed: false
			})
		);
	}),
	rest.delete("/secrets/:token", (req, res, ctx) => {
		return res(
			ctx.status(200),
			ctx.json({
				success: true
			})
		);
	})
);

export default server;
