const defaultConfig = {
	maxWait: "10m",
	resultInterval: "5m",
	timePerPassenger: "10s",
	maxCheckpoints: 5,
};

export async function getNeededCheckpoints(file, config, signal) {
	if (!config)
		config = defaultConfig;

	const formData = new FormData();
	formData.set("file", file);
	formData.set("config", JSON.stringify(config));

	try {
		const resp = await fetch(`${process.env.API_URL}/checkpoints`, {
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
