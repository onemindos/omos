import io
import os
import uuid
import zipfile
import base64
from datetime import datetime, timezone, timedelta

from fastapi import FastAPI, Form, Request
from fastapi.responses import HTMLResponse, StreamingResponse
from fastapi.templating import Jinja2Templates
from cryptography import x509
from cryptography.hazmat.primitives import hashes, serialization
from cryptography.hazmat.primitives.asymmetric import rsa
from cryptography.hazmat.primitives.serialization import pkcs12
from cryptography.x509.oid import NameOID, ExtendedKeyUsageOID
import qrcode

TAK_HOST = os.getenv("TAK_HOST", "tak.tail717ba0.ts.net")
TAK_PORT = os.getenv("TAK_PORT", "8089")
PORTAL_URL = os.getenv("PORTAL_URL", "https://enroll.tail717ba0.ts.net")
CA_CERT_PATH = os.getenv("CA_CERT_PATH", "/certs/ca.pem")
CA_KEY_PATH = os.getenv("CA_KEY_PATH", "/certs/ca.key")
_raw_pass = os.getenv("CA_KEY_PASS", "")
CA_KEY_PASS = _raw_pass.encode() if _raw_pass else None
CERT_PASS = b"atakatak"

app = FastAPI(title="OneMind TAK Portal")
templates = Jinja2Templates(directory="templates")

# In-memory package store: token → zip bytes
# Packages expire on pod restart (acceptable — user rescans if needed)
_packages: dict[str, bytes] = {}


def _load_ca() -> tuple:
    with open(CA_CERT_PATH, "rb") as f:
        ca_cert = x509.load_pem_x509_certificate(f.read())
    with open(CA_KEY_PATH, "rb") as f:
        ca_key = serialization.load_pem_private_key(f.read(), password=CA_KEY_PASS)
    return ca_cert, ca_key


def _issue_cert(callsign: str) -> tuple[bytes, bytes, bytes]:
    """Returns (client_p12, truststore_p12, ca_cert_pem)"""
    ca_cert, ca_key = _load_ca()

    client_key = rsa.generate_private_key(public_exponent=65537, key_size=2048)

    now = datetime.now(timezone.utc)
    subject = x509.Name([
        x509.NameAttribute(NameOID.COUNTRY_NAME, "US"),
        x509.NameAttribute(NameOID.STATE_OR_PROVINCE_NAME, "FL"),
        x509.NameAttribute(NameOID.LOCALITY_NAME, "Tampa"),
        x509.NameAttribute(NameOID.ORGANIZATION_NAME, "OneMindOS"),
        x509.NameAttribute(NameOID.COMMON_NAME, callsign),
    ])

    client_cert = (
        x509.CertificateBuilder()
        .subject_name(subject)
        .issuer_name(ca_cert.subject)
        .public_key(client_key.public_key())
        .serial_number(x509.random_serial_number())
        .not_valid_before(now)
        .not_valid_after(now + timedelta(days=365 * 5))
        .add_extension(x509.BasicConstraints(ca=False, path_length=None), critical=True)
        .add_extension(
            x509.KeyUsage(
                digital_signature=True, key_encipherment=True,
                content_commitment=False, data_encipherment=False,
                key_agreement=False, key_cert_sign=False,
                crl_sign=False, encipher_only=False, decipher_only=False,
            ),
            critical=True,
        )
        .add_extension(
            x509.ExtendedKeyUsage([ExtendedKeyUsageOID.CLIENT_AUTH]),
            critical=True,
        )
        .sign(ca_key, hashes.SHA256())
    )

    client_p12 = pkcs12.serialize_key_and_certificates(
        name=callsign.encode(),
        key=client_key,
        cert=client_cert,
        cas=[ca_cert],
        encryption_algorithm=serialization.BestAvailableEncryption(CERT_PASS),
    )

    truststore_p12 = pkcs12.serialize_key_and_certificates(
        name=b"truststore",
        key=None,
        cert=ca_cert,
        cas=None,
        encryption_algorithm=serialization.BestAvailableEncryption(CERT_PASS),
    )

    ca_pem = ca_cert.public_bytes(serialization.Encoding.PEM)
    return client_p12, truststore_p12, ca_pem


def _build_package(callsign: str) -> str:
    client_p12, truststore_p12, _ = _issue_cert(callsign)

    pkg_uid = str(uuid.uuid4())

    pref_xml = f"""<?xml version='1.0' standalone='yes'?>
<preferences>
  <preference version="1" name="cot_streams">
    <entry key="count" class="class java.lang.Integer">1</entry>
    <entry key="description0" class="class java.lang.String">OneMind OS</entry>
    <entry key="enabled0" class="class java.lang.Boolean">true</entry>
    <entry key="connectString0" class="class java.lang.String">{TAK_HOST},{TAK_PORT},TCP</entry>
  </preference>
  <preference version="1" name="com.atakmap.app_preferences">
    <entry key="displayServerConnectionWidget" class="class java.lang.Boolean">true</entry>
  </preference>
</preferences>"""

    manifest_xml = f"""<MissionPackageManifest version="2">
   <Configuration>
      <Parameter name="uid" value="{pkg_uid}"/>
      <Parameter name="name" value="{callsign}"/>
      <Parameter name="onReceiveImport" value="true"/>
      <Parameter name="onReceiveDelete" value="false"/>
   </Configuration>
   <Contents>
      <Content ignore="false" zipEntry="certs/{callsign}.p12">
         <Parameter name="certPassword" value="atakatak"/>
      </Content>
      <Content ignore="false" zipEntry="certs/truststore-root.p12">
         <Parameter name="certPassword" value="atakatak"/>
      </Content>
      <Content ignore="false" zipEntry="{callsign}.pref"/>
   </Contents>
</MissionPackageManifest>"""

    buf = io.BytesIO()
    with zipfile.ZipFile(buf, "w", zipfile.ZIP_DEFLATED) as zf:
        zf.writestr(f"certs/{callsign}.p12", client_p12)
        zf.writestr("certs/truststore-root.p12", truststore_p12)
        zf.writestr(f"{callsign}.pref", pref_xml)
        zf.writestr("MANIFEST/manifest.xml", manifest_xml)

    token = str(uuid.uuid4())
    _packages[token] = (callsign, buf.getvalue())
    return token


@app.get("/", response_class=HTMLResponse)
async def index(request: Request):
    return templates.TemplateResponse("index.html", {"request": request})


@app.post("/enroll", response_class=HTMLResponse)
async def enroll(request: Request, callsign: str = Form(...)):
    callsign = callsign.strip().replace(" ", "_")[:32]
    if not callsign:
        return templates.TemplateResponse("index.html", {
            "request": request, "error": "Callsign is required."
        })

    token = _build_package(callsign)
    download_url = f"{PORTAL_URL}/download/{token}"

    qr = qrcode.QRCode(error_correction=qrcode.constants.ERROR_CORRECT_M, box_size=8, border=2)
    qr.add_data(download_url)
    qr.make(fit=True)
    img = qr.make_image(fill_color="black", back_color="white")
    buf = io.BytesIO()
    img.save(buf, format="PNG")
    qr_b64 = base64.b64encode(buf.getvalue()).decode()

    return templates.TemplateResponse("enrolled.html", {
        "request": request,
        "callsign": callsign,
        "download_url": download_url,
        "qr_b64": qr_b64,
        "tak_host": TAK_HOST,
        "tak_port": TAK_PORT,
    })


@app.get("/download/{token}")
async def download(token: str):
    entry = _packages.get(token)
    if not entry:
        return HTMLResponse("Package not found or expired. Re-enroll to get a new package.", status_code=404)
    callsign, data = entry
    return StreamingResponse(
        io.BytesIO(data),
        media_type="application/zip",
        headers={"Content-Disposition": f'attachment; filename="{callsign}.zip"'},
    )


@app.get("/health")
async def health():
    return {"status": "ok"}
