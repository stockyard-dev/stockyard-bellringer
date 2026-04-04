package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Bellringer</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--blue:#5b8dd9;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-size:.9rem;letter-spacing:2px}.hdr h1 span{color:var(--rust)}
.main{padding:1.5rem;max-width:960px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(3,1fr);gap:.5rem;margin-bottom:1rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.6rem;text-align:center}
.st-v{font-size:1.2rem;font-weight:700}.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.15rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;align-items:center}
.search{flex:1;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.notif{background:var(--bg2);border:1px solid var(--bg3);padding:.8rem 1rem;margin-bottom:.5rem;transition:border-color .2s}
.notif:hover{border-color:var(--leather)}
.notif-top{display:flex;justify-content:space-between;align-items:flex-start;gap:.5rem}
.notif-title{font-size:.85rem;font-weight:700}
.notif-body{font-size:.7rem;color:var(--cd);margin-top:.2rem}
.notif-meta{font-size:.55rem;color:var(--cm);margin-top:.3rem;display:flex;gap:.6rem;flex-wrap:wrap;align-items:center}
.notif-actions{display:flex;gap:.3rem;flex-shrink:0}
.ch-badge{font-size:.5rem;padding:.12rem .35rem;text-transform:uppercase;letter-spacing:1px;border:1px solid var(--blue);color:var(--blue)}
.stat-pair{display:flex;gap:.15rem;align-items:center;font-size:.55rem}
.stat-pair .v{color:var(--cream)}
.btn{font-size:.6rem;padding:.25rem .5rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:all .2s}
.btn:hover{border-color:var(--leather);color:var(--cream)}.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:460px;max-width:92vw}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.75rem}
</style></head><body>
<div class="hdr"><h1><span>&#9670;</span> BELLRINGER</h1><button class="btn btn-p" onclick="openForm()">+ New Notification</button></div>
<div class="main">
<div class="stats" id="stats"></div>
<div class="toolbar"><input class="search" id="search" placeholder="Search notifications..." oninput="render()"></div>
<div id="list"></div>
</div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()"><div class="modal" id="mdl"></div></div>
<script>
var A='/api',items=[],editId=null;
async function load(){var r=await fetch(A+'/notifications').then(function(r){return r.json()});items=r.notifications||[];renderStats();render();}
function renderStats(){
var total=items.length;
var totalSent=items.reduce(function(s,n){return s+(n.sent_count||0)},0);
var totalClicks=items.reduce(function(s,n){return s+(n.click_count||0)},0);
document.getElementById('stats').innerHTML=[
{l:'Notifications',v:total},{l:'Sent',v:totalSent},{l:'Clicks',v:totalClicks,c:totalClicks>0?'var(--green)':''}
].map(function(x){return '<div class="st"><div class="st-v" style="'+(x.c?'color:'+x.c:'')+'">'+x.v+'</div><div class="st-l">'+x.l+'</div></div>'}).join('');
}
function render(){
var q=(document.getElementById('search').value||'').toLowerCase();
var f=items;
if(q)f=f.filter(function(n){return(n.title||'').toLowerCase().includes(q)||(n.body||'').toLowerCase().includes(q)||(n.channel||'').toLowerCase().includes(q)});
if(!f.length){document.getElementById('list').innerHTML='<div class="empty">No notifications.</div>';return;}
var h='';f.forEach(function(n){
h+='<div class="notif"><div class="notif-top"><div style="flex:1">';
h+='<div class="notif-title">'+esc(n.title)+'</div>';
if(n.body)h+='<div class="notif-body">'+esc(n.body.substring(0,200))+'</div>';
h+='</div><div class="notif-actions">';
h+='<button class="btn btn-sm" onclick="openEdit(''+n.id+'')">Edit</button>';
h+='<button class="btn btn-sm" onclick="del(''+n.id+'')" style="color:var(--red)">&#10005;</button>';
h+='</div></div>';
h+='<div class="notif-meta">';
if(n.channel)h+='<span class="ch-badge">'+esc(n.channel)+'</span>';
h+='<span class="stat-pair">Sent: <span class="v">'+n.sent_count+'</span></span>';
h+='<span class="stat-pair">Clicks: <span class="v">'+n.click_count+'</span></span>';
if(n.url)h+='<span><a href="'+esc(n.url)+'" target="_blank" style="color:var(--blue);font-size:.55rem">'+esc(n.url.substring(0,40))+'</a></span>';
h+='<span>'+ft(n.created_at)+'</span>';
h+='</div></div>';
});
document.getElementById('list').innerHTML=h;
}
async function del(id){if(!confirm('Delete?'))return;await fetch(A+'/notifications/'+id,{method:'DELETE'});load();}
function formHTML(notif){
var i=notif||{title:'',body:'',url:'',icon:'',channel:''};
var isEdit=!!notif;
var h='<h2>'+(isEdit?'EDIT NOTIFICATION':'NEW NOTIFICATION')+'</h2>';
h+='<div class="fr"><label>Title *</label><input id="f-title" value="'+esc(i.title)+'"></div>';
h+='<div class="fr"><label>Body</label><textarea id="f-body" rows="3">'+esc(i.body)+'</textarea></div>';
h+='<div class="row2"><div class="fr"><label>Channel</label><input id="f-channel" value="'+esc(i.channel)+'" placeholder="email, slack, webhook"></div>';
h+='<div class="fr"><label>URL</label><input id="f-url" value="'+esc(i.url)+'" placeholder="https://"></div></div>';
h+='<div class="fr"><label>Icon</label><input id="f-icon" value="'+esc(i.icon)+'"></div>';
h+='<div class="acts"><button class="btn" onclick="closeModal()">Cancel</button><button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Create')+'</button></div>';
return h;
}
function openForm(){editId=null;document.getElementById('mdl').innerHTML=formHTML();document.getElementById('mbg').classList.add('open');}
function openEdit(id){var n=null;for(var j=0;j<items.length;j++){if(items[j].id===id){n=items[j];break;}}if(!n)return;editId=id;document.getElementById('mdl').innerHTML=formHTML(n);document.getElementById('mbg').classList.add('open');}
function closeModal(){document.getElementById('mbg').classList.remove('open');editId=null;}
async function submit(){
var title=document.getElementById('f-title').value.trim();
if(!title){alert('Title is required');return;}
var body={title:title,body:document.getElementById('f-body').value.trim(),url:document.getElementById('f-url').value.trim(),icon:document.getElementById('f-icon').value.trim(),channel:document.getElementById('f-channel').value.trim()};
if(editId){await fetch(A+'/notifications/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
else{await fetch(A+'/notifications',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
closeModal();load();
}
function ft(t){if(!t)return'';try{return new Date(t).toLocaleDateString('en-US',{month:'short',day:'numeric'})}catch(e){return t;}}
function esc(s){if(!s)return'';var d=document.createElement('div');d.textContent=s;return d.innerHTML;}
document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal();});
load();
</script></body></html>`
