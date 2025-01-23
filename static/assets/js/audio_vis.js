
function fillCanvasBar(elt, i, color, data) {
      const max = Math.max(...(data))
      const barwidth = (elt.scrollWidth / data.length) / 2
      const bargap = (elt.scrollWidth / data.length) / 2
      const rect = elt.getBoundingClientRect();
      const ctx = elt.getContext("2d")
      ctx.fillStyle = color;
      const barheight = rect.height * (data[i] / max)
      const x = i * (barwidth + bargap);
      const y = (rect.height / 2) - (barheight) / 2
      ctx.clearRect(x - bargap / 2, 0, barwidth + bargap, rect.height)
      ctx.fillRect(x, y, barwidth, barheight)
  }

  function addVis(canvas) {
    const ctx = canvas.getContext("2d")
    const data = JSON.parse(canvas.getAttribute("audio-vis-data"))
    const barwidth = (canvas.scrollWidth / data.length) / 2
    const bargap = (canvas.scrollWidth / data.length) / 2 

    const dpr = window.devicePixelRatio;
    const rect = canvas.getBoundingClientRect();
    const max = Math.max(...(data))

    // Set the "actual" size of the canvas
    canvas.width = rect.width * dpr;
    canvas.height = rect.height * dpr;

    // Scale the context to ensure correct drawing operations
    ctx.scale(dpr, dpr);

    function fillBar(i, color) {
      ctx.fillStyle = color;
      const barheight = rect.height * (data[i] / max)
      const x = i * (barwidth + bargap);
      const y = (rect.height / 2) - (barheight) / 2
      ctx.clearRect(x - bargap / 2, 0, barwidth + bargap, rect.height)
      ctx.fillRect(x, y, barwidth, barheight)
    }

    function fillBars(color, divisor) {
      ctx.fillStyle = color
      for (let i = 0; i < data.length; i++) {
        barheight = rect.height * (data[i] / max) / divisor
        ctx.fillRect(i * (barwidth + bargap), (rect.height / 2) - barheight / 2, barwidth, barheight)
      }
    }
    fillBars("white", 1);

    const statusArr = new Array(data.length).fill(false)
    
  canvas.timeUpdate = () => {
      const progress = audio.currentTime / audio.duration
      const index = data.length * progress

      for (let i = 0; i < statusArr.length; i++) {
        if (statusArr[i] == false && i <= index) {
          fillBar(i, "#898b8f")
          statusArr[i] = true;
        } else if (statusArr[i] == true && i > index) {
          fillBar(i, "white")
          statusArr[i] = false;
        }
      }
    }

  canvas.onEnd = () => {
      for (let i = 0; i < data.length; i++) {
        fillBar(i, "white")
      }
    }
  }


// htmx.on('htmx:afterProcessNode', 

    // canvas.addEventListener('mousemove', function (e) {
    //   canvasX = e.offsetX
    //   canvasY = e.offsetY
    //   const index = Math.floor(canvasX / (barwidth + bargap))
    //
    //   for (let i = 0; i < statusArr.length; i++) {
    //     if (statusArr[i] == false && i <= index) {
    //       fillBar(i, "darkgreen")
    //       statusArr[i] = true;
    //     } else if (statusArr[i] == true && i > index) {
    //       fillBar(i, "limegreen")
    //       statusArr[i] = false;
    //     }
    //   }
    // })


/*
  window.addEventListener("popstate", (e) => {
    const canvases = document.querySelectorAll(".audio-vis");
    canvases.forEach((el) => {
      addVis(el);
      el.visLoaded = true;
    })
  })

  window.onload = () => {
    const canvases = document.querySelectorAll(".audio-vis");
    canvases.forEach((el) => {
      addVis(el);
      el.visLoaded = true;
    })
  }

  window.addEventListener("htmx:historyRestore", (e) => {
    // console.log("history restore event")
    const canvases = document.querySelectorAll(".audio-vis");
    canvases.forEach((el) => {
      addVis(el);
      el.visLoaded = true;
    })
  })
  window.addEventListener("htmx:afterSettle", (e) => {
    const newElt = e.detail.elt
    newCanvases = newElt.querySelectorAll("canvas.audio-vis")
    newCanvases.forEach((el) => {
      if (!el.visLoaded) {
        addVis(el);
      }
      el.visLoaded = true
    })
  })

*/
