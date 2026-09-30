"""OneMind TAK-ML Vision Worker — pulls RTSP, runs YOLOv5s, publishes to NATS."""
import asyncio, json, os, time, logging, ssl
from datetime import datetime, timezone
import cv2
import numpy as np
import torch
import nats

logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s")
logger = logging.getLogger("vision-worker")

NATS_URL = os.environ.get("NATS_URL", "nats://nats.nats.svc.cluster.local:4222")
NATS_USER = os.environ.get("NATS_USER", "onemind")
NATS_PASS = os.environ.get("NATS_PASSWORD", "")
NATS_CA = os.environ.get("NATS_CA", "/etc/nats-ca/ca.crt")
MODEL_PATH = os.environ.get("MODEL_PATH", "/models/visdrone.torchscript")
CONFIDENCE = float(os.environ.get("CONFIDENCE_THRESHOLD", "0.45"))
INTERVAL = float(os.environ.get("SAMPLE_INTERVAL", "2.0"))
IMG_SIZE = int(os.environ.get("IMG_SIZE", "640"))
LABELS = ["pedestrian","people","bicycle","car","van","truck","tricycle","awning-tricycle","bus","motor"]
FEEDS = json.loads(os.environ.get("FEEDS", '{}'))

def xywh2xyxy(x):
    y = np.copy(x)
    y[:,0]=x[:,0]-x[:,2]/2; y[:,1]=x[:,1]-x[:,3]/2
    y[:,2]=x[:,0]+x[:,2]/2; y[:,3]=x[:,1]+x[:,3]/2
    return y

def nms_boxes(boxes, scores, iou_thr=0.45):
    x1,y1,x2,y2 = boxes[:,0],boxes[:,1],boxes[:,2],boxes[:,3]
    areas = (x2-x1)*(y2-y1)
    order = scores.argsort()[::-1]
    keep = []
    while order.size > 0:
        i = order[0]; keep.append(i)
        if order.size == 1: break
        xx1=np.maximum(x1[i],x1[order[1:]]); yy1=np.maximum(y1[i],y1[order[1:]])
        xx2=np.minimum(x2[i],x2[order[1:]]); yy2=np.minimum(y2[i],y2[order[1:]])
        inter = np.maximum(0,xx2-xx1)*np.maximum(0,yy2-yy1)
        iou = inter/(areas[i]+areas[order[1:]]-inter)
        order = order[np.where(iou<=iou_thr)[0]+1]
    return np.array(keep)

def parse_detections(output, conf_thr=0.45):
    if output.ndim == 3: output = output[0]
    mask = output[:,4] > conf_thr
    output = output[mask]
    if len(output)==0: return []
    scores = output[:,5:]*output[:,4:5]
    cls_ids = np.argmax(scores, axis=1)
    confs = scores[np.arange(len(cls_ids)), cls_ids]
    mask2 = confs > conf_thr
    output,cls_ids,confs = output[mask2],cls_ids[mask2],confs[mask2]
    if len(output)==0: return []
    boxes = xywh2xyxy(output[:,:4])
    dets = []
    for cid in np.unique(cls_ids):
        m = cls_ids==cid
        keep = nms_boxes(boxes[m], confs[m])
        for idx in keep:
            dets.append({"class":LABELS[cid] if cid<len(LABELS) else str(cid),
                "confidence":float(confs[m][idx]),"bbox":boxes[m][idx].tolist()})
    dets.sort(key=lambda x:x["confidence"], reverse=True)
    return dets[:50]

async def process_feed(nc, model, feed_id, url):
    logger.info(f"Starting feed: {feed_id}")
    while True:
        cap = cv2.VideoCapture(url)
        if not cap.isOpened():
            logger.warning(f"[{feed_id}] Cannot open {url}, retry 10s")
            await asyncio.sleep(10); continue
        logger.info(f"[{feed_id}] Connected")
        fails = 0
        while True:
            ret, frame = cap.read()
            if not ret:
                fails += 1
                if fails > 10: break
                await asyncio.sleep(1); continue
            fails = 0; t0 = time.time()
            img = cv2.resize(frame, (IMG_SIZE, IMG_SIZE))
            img = cv2.cvtColor(img, cv2.COLOR_BGR2RGB)
            tensor = torch.from_numpy(img).permute(2,0,1).float()/255.0
            tensor = tensor.unsqueeze(0)
            with torch.no_grad():
                out = model(tensor)
            if isinstance(out, tuple): out = out[0]
            dets = parse_detections(out.numpy(), CONFIDENCE)
            ms = (time.time()-t0)*1000
            if dets:
                await nc.publish(f"fabric.vision.{feed_id}.detect",
                    json.dumps({"timestamp":datetime.now(timezone.utc).isoformat(),
                        "feed_id":feed_id,"detections":dets,"inference_ms":round(ms,1)}).encode())
                logger.info(f"[{feed_id}] {len(dets)} dets in {ms:.0f}ms")
            await asyncio.sleep(INTERVAL)
        cap.release(); await asyncio.sleep(5)

async def main():
    logger.info("OneMind Vision Worker starting")
    if not os.path.exists(MODEL_PATH):
        logger.error(f"Model not found: {MODEL_PATH}"); return
    model = torch.jit.load(MODEL_PATH, map_location="cpu"); model.eval()
    dummy = torch.zeros(1,3,IMG_SIZE,IMG_SIZE)
    with torch.no_grad(): model(dummy)
    logger.info("Model loaded")
    opts = {"servers":[NATS_URL]}
    if NATS_USER: opts["user"]=NATS_USER
    if NATS_PASS: opts["password"]=NATS_PASS
    if os.path.exists(NATS_CA):
        ctx = ssl.create_default_context(cafile=NATS_CA); opts["tls"]=ctx
    nc = await nats.connect(**opts)
    logger.info("NATS connected")
    tasks = [process_feed(nc, model, fid, url) for fid, url in FEEDS.items()]
    await asyncio.gather(*tasks)

if __name__ == "__main__":
    asyncio.run(main())
