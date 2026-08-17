FROM python:3.10-slim

WORKDIR /app

RUN apt-get update && \
    apt-get install -y ffmpeg build-essential && apt-get clean

COPY requirements.txt ./
RUN pip install --no-cache-dir --upgrade pip && pip install -r requirements.txt --no-cache-dir --extra-index-url \
    https://download.pytorch.org/whl/cpu

COPY . .

EXPOSE 5000

CMD ["flask", "run", "--host=0.0.0.0"]
