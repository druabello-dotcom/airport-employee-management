import chokidar from "chokidar";

const clients = new Set();

chokidar.watch("./src", { usePolling: true }).on("change", path => {
	console.log(`[reload] ${path} changed`);
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

		const filePath =
			url.pathname === "/" ? "./src/index.html" : `./src${url.pathname}`;

		if (filePath.endsWith(".html")) {
			let html = await Bun.file(filePath).text();
			html = html.replace(
				"</body>",
				`<script>
				const es = new EventSource("/__reload");
				es.onmessage = () => location.reload();
			</script></body>`
			);
			return new Response(html, { headers: { "Content-Type": "text/html" } });
		}

		return new Response(Bun.file(filePath));
	},
});

console.log("Dev server running at http://localhost:3000");
