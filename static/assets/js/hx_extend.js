htmx.defineExtension('audio', {
  onEvent: function(name, event) {
    console.log(`name: ${name}\nevent: ${string(event)}`)
  }
})
