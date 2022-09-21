async function downloadJson(url) {
    const file = await fetch(url);
    return await file.json();
}

async function runWasm(fileUrl) {
    try {
        const file = await fetch(fileUrl);
        const buffer = await file.arrayBuffer();
        const go = new Go();
        const { instance } = await WebAssembly.instantiate(buffer, go.importObject);
        go.run(instance);
    } catch (err) {
        console.error(err);
    }
}

