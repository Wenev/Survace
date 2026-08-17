from flask import Flask, request, jsonify
import whisper
import ffmpeg
import os
import tempfile
from flask_cors import CORS

app = Flask(__name__)
CORS(app)
model = whisper.load_model("small")

@app.route('/caption', methods=['POST'])
def caption():
    target_lang = request.args.get('target_lang', 'original')
    if 'file' not in request.files:
        return jsonify({'error': 'No file uploaded'}), 400

    file = request.files['file']
    if file.filename == '':
        return jsonify({'error': 'No selected file'}), 400

    if not file.filename.lower().endswith('.mp4'):
        return jsonify({'error': 'Only MP4 files are supported'}), 400

    with tempfile.NamedTemporaryFile(delete=False, suffix='.mp4') as temp_video:
        file.save(temp_video.name)
        wav_path = temp_video.name + '.wav'

        try:
            try:
                probe = ffmpeg.probe(temp_video.name)
                audio_streams = [stream for stream in probe['streams'] if stream['codec_type'] == 'audio']
                if not audio_streams:
                    return jsonify({'error': 'No audio stream found in video'}), 400

                (
                    ffmpeg.input(temp_video.name)
                    .output(wav_path, ac=1, ar=16000, acodec='pcm_s16le')
                    .run(overwrite_output=True, quiet=True, capture_stderr=True)
                )
            except ffmpeg.Error as e:
                return jsonify({'error': f'FFmpeg error: {e.stderr.decode()}'}), 400

            if not os.path.exists(wav_path) or os.path.getsize(wav_path) == 0:
                return jsonify({'error': 'Audio extraction failed'}), 400

            audio_info = ffmpeg.probe(wav_path)
            duration = float(audio_info['format']['duration'])
            if duration < 0.5:
                return jsonify({'error': 'Audio is too short (needs at least 0.5 seconds)'}), 400

            try:
                if target_lang == 'en':
                    result = model.transcribe(wav_path, task="translate", language='en')
                elif target_lang == 'es':
                    result = model.transcribe(wav_path, task="translate", language='es')
                else:
                    result = model.transcribe(wav_path)

                segments = [
                    {
                        'start': segment['start'],
                        'end': segment['end'],
                        'text': segment['text']
                    }
                    for segment in result.get('segments', [])
                ]
                text = result.get('text', '')

            except Exception as e:
                return jsonify({'error': f'Transcription failed: {str(e)}'}), 500

        finally:
            for path in [temp_video.name, wav_path]:
                if os.path.exists(path):
                    try:
                        os.remove(path)
                    except:
                        pass

    return jsonify({'segments': segments, 'text': text})

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
