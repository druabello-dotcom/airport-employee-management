export async function getNeededCheckpoints(file, config, signal) {
	const formData = new FormData();
	formData.set("file", file);
	formData.set("config", config);

	try {
		const resp = await fetch(`${Bun.env.API_URL}/checkpoints`, {
			method: "POST",
			body: formData,
			signal: signal,
		})
		if (!resp.ok)
			throw new Error(`checkpoints request failed: ${await resp.text()}`);
		const data = await resp.json();
		return data;
	} catch (err) {
		if (err.name === "AbortError") return;
		throw err;
	}
}
