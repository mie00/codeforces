function zip(a, b) {
    return a.map(function(e, i) {
        return { "input": e, "output": b[i] };
    });
}
var inp = Array.from(document.getElementsByClassName("input")).map(x => (x.tagName === "PRE" ? x : x.getElementsByTagName("pre")[0]).innerText.trim())
var out = Array.from(document.getElementsByClassName("answer")).map(x => (x.tagName === "PRE" ? x : x.getElementsByTagName("pre")[0]).innerText.trim())
if (!out.length) {
    out = Array.from(document.getElementsByClassName("output")).map(x => (x.tagName === "PRE" ? x : x.getElementsByTagName("pre")[0]).innerText.trim())
}
var json = JSON.stringify(zip(inp, out).filter(x => x.input != '?' || x.output != '?'))

function fallbackCopyTextToClipboard(text) {
    var textArea = document.createElement("textarea");
    textArea.value = text;
    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();

    try {
        var successful = document.execCommand('copy');
        var msg = successful ? 'successful' : 'unsuccessful';
        console.log('Fallback: Copying text command was ' + msg);
    } catch (err) {
        console.error('Fallback: Oops, unable to copy', err);
    }

    document.body.removeChild(textArea);
}

function copyTextToClipboard(text) {
    if (!navigator.clipboard) {
        fallbackCopyTextToClipboard(text);
        return;
    }
    navigator.clipboard.writeText(text).then(function() {
        console.log('Async: Copying to clipboard was successful!');
    }, function(err) {
        console.error('Async: Could not copy text: ', err);
    });
}

copyTextToClipboard(json);