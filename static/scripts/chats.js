function switchFrame(frameId) {
    const frames = document.querySelectorAll(".frame");
    frames.forEach(frame => {
        frame.classList.add("hidden")
    });
    document.querySelector("#" + frameId).classList.remove("hidden");
    return
}

// UNUSED
(() => {
    const withSwitchFrame = document.querySelectorAll("[gm-setframe]")
    withSwitchFrame.forEach(e => {
        const frames = document.querySelectorAll(".frame")
        const eventType = e.getAttribute("gm-trigger") || "onclick";
        e.addEventListener(eventType, () => {
            const frameId = e.getAttribute("gm-setframe");
            if (!frameId) {
                return;
            }

            frames.forEach(f => {
                f.classList.add("hidden");
            })

            document.querySelector("#" + frameId).classList.remove("hidden");
        })
    })
})
