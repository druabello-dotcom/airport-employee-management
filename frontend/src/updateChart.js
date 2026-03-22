let chartExist = null;

export function updateChart(dataObject) {
	console.log("update chart function called", dataObject);
	const times = [];
	const checkpoints = [];
	for (let i = 0; i < dataObject.length; i++) {
		let totalMinutes = dataObject[i].time / (60 * Math.pow(10, 9));
		let hour = Math.floor(totalMinutes / 60);
		let minutes = totalMinutes % 60;
		if (minutes < 10) 
			times.push(`${hour}:0${minutes}`);
		else 
			times.push(`${hour}:${minutes}`);
		checkpoints.push(dataObject[i].minOpen);
	}

	if (chartExist) chartExist.destroy();
	chartExist  = new Chart("myChart", {
		type: "line",
		data: {
			labels: times,
			datasets: [
				{
					label: "Employee Managment",
					borderColor: "#2596be",
					data: checkpoints,
				}
			]
		},
		options: {
			responsive: true,
			maintainAspectRatio: false,
			plugins: {
				legend: {
					display: false
				}
			}
		}
	})
}