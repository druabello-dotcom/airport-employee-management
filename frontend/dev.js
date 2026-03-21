import chokidar from "chokidar";

const clients = new Set();

async function buildFrontend() {
	const result = await Bun.build({
		entrypoints: ["./src/main.js"],
		outdir: "./dist",
		naming: "main.js",
		define: {
			"process.env.API_URL": JSON.stringify(Bun.env.API_URL),
		},
	});

	if (!result.success) {
		console.error("Build failed", result.logs);
	}
}

// Initial build.
await buildFrontend();

chokidar.watch("./src", { usePolling: true }).on("change", async path => {
	console.log(`[reload] ${path} changed`);
	await buildFrontend();

	for (const client of clients) {
		try {
			client.enqueue("data: reload\n\n");
		} catch {
			clients.delete(client);
		}
	}
});

Bun.serve({
	port: 3000,
	async fetch(req) {
		const url = new URL(req.url);

		if (url.pathname === "/__reload") {
			let controller;
			const stream = new ReadableStream({
				start(c) {
					controller = c;
					clients.add(controller);
					req.signal.addEventListener("abort", () => {
						clients.delete(controller);
					});
				},
			});
			return new Response(stream, {
				headers: {
					"Content-Type": "text/event-stream",
					"Cache-Control": "no-cache",
					Connection: "keep-alive",
				},
			});
		}

		if (url.pathname === "/") {
			let html = await Bun.file("./src/index.html").text();
			html = html.replace(
				"</body>",
				`<script>
				const es = new EventSource("/__reload");
				es.onmessage = () => location.reload();
			</script></body>`
			);
			return new Response(html, { headers: { "Content-Type": "text/html" } });
		}

		if (url.pathname === "/main.js")
			return new Response(Bun.file("./dist/main.js"));

		return new Response(Bun.file(`./src${url.pathname}`));
	},
});

console.log("Dev server running at http://localhost:3000");
