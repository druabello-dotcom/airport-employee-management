document.getElementById('form').addEventListener('submit', checkpointReq);
async function checkpointReq(file) {
	if (!file) {
		console.warn('No file selected');
		return;
	}

	const formData = new FormData();
	formData.append('user-file', file);
}