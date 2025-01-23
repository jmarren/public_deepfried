onmessage = (e) => {
  const file = e.data

  const fileToUpload = new FileToUpload(file.file, file.type, 0)
  fileToUpload.upload()
}

class Upload {
  constructor(filesToUpload) {
    this.filesToUpload = filesToUpload
  }
  uploadFiles() {
    for (let i = 0; i < this.filesToUpload.length; i++) {
      this.filesToUpload[i].upload()
    }
  }
}

class FileToUpload {
  constructor(file, type, id) {
    this.file = file
    this.UploadType = type
    this.fileName = file.name
    this.MAX_READ_LENGTH = 1024 * 512
    this.writeStream = this.createWriteStream()
    this.chunkIndex = 0
    this.id = id
    this.progress = 0
  }
  incrementChunkIndex() {
    this.chunkIndex++
  }

  async upload() {
        this.file.stream().pipeTo(this.writeStream)
  }
  async postChunks(fullChunk) {
    const maxChunkSize = this.MAX_READ_LENGTH
    const numChunks = Math.ceil(fullChunk.byteLength / maxChunkSize);

    for (let i = 0; i < numChunks; i++) {
        const start = i * maxChunkSize
        let end = start + maxChunkSize
        if (end > fullChunk.byteLength) {
          end = fullChunk.byteLength
        }
        const chunk = fullChunk.slice(start, end)
        await this.postChunk(chunk)
    }


  }
  async postChunk(chunk) {
    const index = `${this.chunkIndex}`
    const fileId = `${this.id}`
    const fileName = `${this.fileName}`
    const fileSize = `${this.file.size}`
    const fileType = `${this.file.type}`
    let message 
    try {

      const chunkSize = `${chunk.byteLength}`
      const url = `/upload/append/${this.UploadType}`

      const response = await fetch(url, {
        method: 'POST',
        body: chunk,
        headers: {
          'Content-Type': fileType,
          'X-index': index,
          'X-total-size': fileSize,
          'X-file-name': fileName,
        }
      })
      if (!response.ok) {
        throw new Error("failed to upload chunk");
      }
      const res = await response.text()
      message = res

    } catch (e) {
      console.error(e)
    } finally {
      this.chunkIndex++
      this.progress = this.progress + chunk.byteLength / fileSize;
      postMessage([this.fileName, this.progress, message])
    }
  }

  createWriteStream() {
    const queuingStrategy = new CountQueuingStrategy({ highWaterMark: 1})

    const underlyingSink = {
      write: (chunk) => {
        return this.postChunks(chunk)
      },
      close: () => {
        postMessage([this.fileName, this.progress, "complete"])
      },
      abort: (err) => {
        console.error(err)
      },
    }
    const writableStream = new WritableStream(underlyingSink, queuingStrategy)
    return writableStream
  }
}
