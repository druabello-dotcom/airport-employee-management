let chartExist = null;

export function updateChart(dataObject) {
	const times = [];
	const checkpoints = [];
	const waitTimes = [];

	dataObject.forEach(data => {
		const totalMinutes = data.time / (60 * Math.pow(10, 9));
		const hour = Math.floor(totalMinutes / 60);
		const minutes = totalMinutes % 60;
		if (minutes < 10)
			times.push(`${hour}:0${minutes}`);
		else
			times.push(`${hour}:${minutes}`);

		checkpoints.push(data.minOpen);

		const waitMinutes = data.wait / (60 * Math.pow(10, 9));
		waitTimes.push(waitMinutes);
	});

	if (chartExist) chartExist.destroy();
	chartExist = new Chart("myChart", {
		type: "line",
		data: {
			labels: times,
			datasets: [
				{
					label: "Checkpoints",
					borderColor: "#2596be",
					data: checkpoints,
					yAxisID: "y1",
				}, {
					label: "Wait Time",
					data: waitTimes,
					yAxisID: "y2",
				},
			]
		},
		options: {
			responsive: true,
			maintainAspectRatio: false,
			plugins: {
				legend: {
					display: false
				}
			},
		scales: {
			y1: {
				type: "linear",
				position: "left",
			},
			y2: {
				type: "linear",
				position: "right",
				grid: {
					drawOnChartArea: false,
				},
			},
		},
		}
	})
}
