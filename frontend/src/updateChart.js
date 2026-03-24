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
					data: checkpoints,
					yAxisID: "y1",
					tension: 0.2,
				}, {
					label: "Wait Time",
					data: waitTimes,
					yAxisID: "y2",
					tension: 0.2,
				},
			]
		},
		options: {
			responsive: true,
			maintainAspectRatio: false,
			scales: {
				y1: {
					type: "linear",
					position: "left",
					min: 0,
				},
				y2: {
					type: "linear",
					position: "right",
					min: 0,
					grid: {
						drawOnChartArea: false,
					},
					ticks: {
						callback: value => {
							return value + " min";
						},
					},
				},
			},
			plugins: {
				tooltip: {
					callbacks: {
						label: ctx => {
							const prefix = ctx.dataset.label + ": ";
							let val = ctx.parsed.y;
							let suffix = "";

							if (ctx.dataset.yAxisID === "y2") {
								val = Math.round(ctx.parsed.y * 10) / 10;
								suffix = " min";
							}

							return prefix + val + suffix;
						},
					}
				}
			}
		},
	})
}
