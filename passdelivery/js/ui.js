/* PassDelivery — общие мелочи интерфейса */
export const $ = (id) => document.getElementById(id);

export function say(id, text, kind){ const el = $(id); el.textContent = text || ''; el.className = 'msg' + (kind ? ' ' + kind : ''); }
export function busy(btn, on){ btn.disabled = on; }

// Для 400 бэк кладёт текст валидации в error, а в message — общий заголовок
export function errText(e){ return e.status === 400 && e.body?.error ? e.body.error : e.message; }

export function initials(name){
  return (name || '').trim().split(/\s+/).slice(0, 2).map(w => w[0] || '').join('').toUpperCase() || '·';
}
