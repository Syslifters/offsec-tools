# shareranger/utils/stream_writer.py
import json, csv, threading, queue, os
from typing import Optional, Dict, Any, IO
from shareranger.logging.logger import log_error

Event = Dict[str, Any]


class StreamResultWriter:
    """
    Many producers enqueue events; a single consumer thread writes them.
    Files are owned by this class: opened in start(), closed in stop().
    """

    def __init__(self, ndjson_path: Optional[str], csv_path: Optional[str]):
        self.ndjson_path = ndjson_path
        self.csv_path = csv_path

        self.q: "queue.Queue[Event]" = queue.Queue()
        self._stop = threading.Event()
        self._thr = threading.Thread(target=self._run, daemon=True)

        self._ndjson_file: Optional[IO] = None
        self._csv_file: Optional[IO] = None
        self._csv_writer: Optional[csv.DictWriter] = None

    # ------------ lifecycle ------------
    def start(self):
        if self.ndjson_path:
            self._ndjson_file = open(self.ndjson_path, "a", encoding="utf-8")

        if self.csv_path:
            is_new = not os.path.exists(self.csv_path)
            self._csv_file = open(self.csv_path, "a", newline="", encoding="utf-8")
            self._csv_writer = csv.DictWriter(
                self._csv_file,
                fieldnames=[
                    "severity",
                    "type",
                    "host",
                    "share",
                    "path",
                    "id",
                    "location",
                    "pattern",
                    "snippet",
                ],
            )
            if is_new:
                self._csv_writer.writeheader()
                self._csv_file.flush()

        self._thr.start()

    def stop(self):
        # Ask consumer to drain and exit
        self._stop.set()
        self.q.put({"t": "__flush__"})
        self._thr.join()

        # Close files we own
        try:
            if self._ndjson_file:
                self._ndjson_file.flush()
                self._ndjson_file.close()
        finally:
            self._ndjson_file = None

        try:
            if self._csv_file:
                self._csv_file.flush()
                self._csv_file.close()
        finally:
            self._csv_file = None
            self._csv_writer = None

    # ------------ enqueue helpers ------------
    def enqueue_host(self, host):
        self.q.put(
            {"t": "host", "name": host.name, "online": getattr(host, "online", None)}
        )

    def enqueue_share(self, host_name: str, share):
        self.q.put(
            {
                "t": "share",
                "host": host_name,
                "name": share.name,
                "unc_path": share.unc_path,
                "accessible": getattr(share, "accessible", None),
            }
        )

    def enqueue_match(
        self, host_name: str, share_name: str, match_dict: Dict[str, Any]
    ):
        self.q.put(
            {"t": "match", "host": host_name, "share": share_name, "match": match_dict}
        )

    def enqueue_snapshot(self, hosts_list_dicts):
        self.q.put({"t": "snapshot", "hosts": hosts_list_dicts})

    # ------------ consumer thread ------------
    def _run(self):
        try:
            while True:
                ev = self.q.get()
                if ev.get("t") == "__flush__":
                    if self._stop.is_set():
                        break
                    continue

                # NDJSON: write all non-snapshot events (you can include snapshots if desired)
                if self._ndjson_file and ev["t"] != "snapshot":
                    self._ndjson_file.write(json.dumps(ev, ensure_ascii=False) + "\n")
                    self._ndjson_file.flush()

                # CSV: only matches
                if self._csv_writer and ev["t"] == "match":
                    m = ev["match"]
                    self._csv_writer.writerow(
                        {
                            "severity": m.get("severity", ""),
                            "type": m.get("type", ""),
                            "host": ev["host"],
                            "share": ev["share"],
                            "path": m.get("unc_path", ""),
                            "id": m.get("id", ""),
                            "location": m.get("location", ""),
                            "pattern": m.get("pattern", ""),
                            "snippet": m.get("snippet", ""),
                        }
                    )
                    self._csv_file.flush()
        except Exception as e:
            # Do not close files here; stop() will handle cleanup.
            log_error(f"StreamResultWriter crashed: {e}")
