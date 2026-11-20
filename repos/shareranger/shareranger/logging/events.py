# shareranger/events.py
from typing import Optional, Dict, Any
from shareranger.logging.stream_writer import StreamResultWriter

_writer: Optional[StreamResultWriter] = None


def start(ndjson_path: str, csv_path: str) -> StreamResultWriter:
    """Create + start the singleton writer."""
    global _writer
    if _writer is None:
        _writer = StreamResultWriter(ndjson_path, csv_path)
        _writer.start()
    return _writer


def stop() -> None:
    """Stop and clear the singleton writer."""
    global _writer
    if _writer is not None:
        _writer.stop()
        _writer = None


# ---- convenience emitters (safe no-ops if not started) ----
def emit_host(host) -> None:
    if _writer:
        _writer.enqueue_host(host)


def emit_share(host_name: str, share) -> None:
    if _writer:
        _writer.enqueue_share(host_name, share)


def emit_match(host_name: str, share_name: str, match_dict: Dict[str, Any]) -> None:
    if _writer:
        _writer.enqueue_match(host_name, share_name, match_dict)


def emit_match_for_share(share, match_dict: Dict[str, Any]) -> None:
    """Uses share._host_name if available."""
    host_name = getattr(share, "_host_name", None)
    if host_name and _writer:
        _writer.enqueue_match(host_name, share.name, match_dict)
