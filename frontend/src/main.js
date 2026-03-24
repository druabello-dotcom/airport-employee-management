import { defaultConfig, getNeededCheckpoints } from "./api.js";
import { updateChart } from "./updateChart.js";

document.getElementById("uploadTrigger").addEventListener("click", () => {
	document.getElementById("file").click();
})

document.getElementById("file").addEventListener("change", (event) => {
    const file = event.target.files[0];
    const fileStatus = document.getElementById("file-name");
    if (file) {
        fileStatus.textContent = file.name;
    } else {
        fileStatus.textContent = "No file selected";
    }
});

document.getElementById("form").addEventListener("submit", async event => {
	event.preventDefault();

	const file = document.getElementById("file").files[0];
	if (!file) {
		console.warn("No file selected");
		return;
	}
	let results = await getNeededCheckpoints(file);
	updateChart(results);
});

const displayMaxWait = document.getElementById("max-wait-time-value");
displayMaxWait.textContent = defaultConfig.maxWait;

const displayResultInterval = document.getElementById("result-interval-value");
displayResultInterval.textContent = defaultConfig.resultInterval;

const displayTimePerPassenger = document.getElementById("time-per-passenger-value");
displayTimePerPassenger.textContent = defaultConfig.timePerPassenger;

const displayMaxCheckpoints = document.getElementById("max-checkpoints-value");
displayMaxCheckpoints.textContent = defaultConfig.maxCheckpoints;

// input default config object in api.js
document.getElementById("submitMaxCheckpoints").onclick = function() {
	let input = document.getElementById("input-max-checkpoints").value;
	defaultConfig.maxCheckpoints = Number(input);
	displayMaxCheckpoints.textContent = defaultConfig.maxCheckpoints;
}
document.getElementById("submit-time-per-passenger").onclick = function() {
	let input = document.getElementById("input-time-per-passenger").value;
	defaultConfig.timePerPassenger = `${Number(input)}s`;
	displayTimePerPassenger.textContent = defaultConfig.timePerPassenger;
}
document.getElementById("submit-max-wait-time").onclick = function() {
	let input = document.getElementById("input-max-wait-time").value;
	defaultConfig.maxWait = `${Number(input)}m`;
	displayMaxWait.textContent = defaultConfig.maxWait;
}
document.getElementById("submit-result-interval").onclick = function() {
	let input = document.getElementById("input-result-interval").value;
	defaultConfig.resultInterval = `${Number(input)}m`;
	displayResultInterval.textContent = defaultConfig.resultInterval;
}