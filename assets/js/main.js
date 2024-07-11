
// ================================================================ //
// ========== PLATAFORMA ========================================== //
function getPlatform() {
    // Modern way of detecting
    if (typeof navigator.userAgentData !== 'undefined' && navigator.userAgentData != null) {
        return navigator.userAgentData.platform;
    }
    // Deprecated fallback
    if (typeof navigator.platform !== 'undefined') {
        if (typeof navigator.userAgent !== 'undefined' && /android/.test(navigator.userAgent.toLowerCase())) {
            // android device’s navigator.platform is often set as 'linux', so let’s use userAgent for them
            return 'android';
        }
        return navigator.platform;
    }
    return 'unknown';
}
// let platform = getPlatform().toLowerCase();
let esMac = /mac/.test(getPlatform().toLowerCase()); // Mac desktop
// let esIOS = ['iphone', 'ipad', 'ipod'].indexOf(platform) >= 0; // Mac iOs
// let esApple = esMac || esIOS; // Apple device (desktop or iOS)
// let esWindows = /win/.test(platform); // Windows
// let esAndroid = /android/.test(platform); // Android
// let esLinux = /linux/.test(platform); // Linux



// ================================================================ //
// ========== SHORTCUTS =========================================== //

// Guardar con "CTRL + S" y hx-trigger="submit,cmdGuardar"
// Se debe prevenir default en keydown pero lanzar el evento
// en keyup para solo lanzarlo una vez.
document.addEventListener("keydown", function(e) {
	if ((esMac ? e.metaKey : e.ctrlKey) && e.code === 'KeyS') {
		e.preventDefault();
	}
}, false);
document.addEventListener("keyup", function(e) {
	if ((esMac ? e.metaKey : e.ctrlKey) && e.code === 'KeyS' ) {
		e.target.dispatchEvent(new Event("cmdGuardar", { bubbles: true }))
	}
}, false);


// ================================================================ //
// ================================================================ //


function delay(time) {
	return new Promise(resolve => setTimeout(resolve, time));
}

function quitar(element) {
	console.log("quitando")
	setTimeout(() => {
		console.log("bye")
		element.remove();
	}, 200);
}

window.onload = () => {
	// console.log("page is fully loaded");
	
	htmx.on("htmx:responseError", function(event) {
		alert(`Error ${event.detail.xhr.response}`)
	});
	htmx.on("htmx:sendError", function(event) {
		alert(`Error de red: no se puede conectar con el servidor.`)
	});
	htmx.on("htmx:timeout", function(event) {
		alert(`Error: se agotó el tiempo de espera para la respuesta del servidor.`)
	});
};


/**
 * 
 * @param {string} str string para quitar acentos, espacios, trim, diacríticos.
 * @returns string
 */
// normalizar quita acentos, espacios adicionales y mayúsculas.
function normalizar(str) {
	return str.normalize('NFD').replace(/\p{Diacritic}/gu, '').toLowerCase().trim().replace(/\s+/g, ' ')
}


// Filtrar entidades por nombre
// 
// Uso:
// <input type="search" class="" oninput="filtrarArticles(this.value)" placeholder="Buscar entidad...">
//
function filtrarArticles(qryText) {
	let cards = document.getElementsByTagName("article")
	if (cards.length < 1) {
		console.log("no hay article para filtrar")
		return
	}
	qryText = normalizar(qryText)
	for (let card of cards) {
		if ( qryText.length < 2 ) {
			card.classList.remove("hidden")
			continue
		}
		let cardText = normalizar(card.getElementsByTagName("h2")[0].textContent)
		if ( cardText.includes(qryText) ) {
			card.classList.remove("hidden")
		} else {
			card.classList.add("hidden")
		}
	}
}


function filtrarRows(qryText, tableID) {
	let rows = document.getElementById(tableID).getElementsByTagName("tbody")[0].getElementsByTagName("tr")
	if (rows.length < 1) {
		return
	}
	qryText = normalizar(qryText)
	for (let row of rows) {
		if ( qryText.length < 2 ) {
			row.classList.remove("hidden")
			continue
		}
		let rowText = normalizar(row.textContent)
		if ( rowText.includes(qryText) ) {
			row.classList.remove("hidden")
		} else {
			row.classList.add("hidden")
		}
	}
}

// Mostar checkboxes para seleccionar varias entidades.
function showCheckboxs() {
	let chbxs = document.getElementsByClassName("filtro_chkbox")
	for (let cb of chbxs) {
		cb.classList.remove("hidden")
	}
	document.getElementById("showCheckboxsBtn").classList.add("hidden")
	document.getElementById("applyCheckboxsBtn").classList.remove("hidden")
}


/**
 * Ordena una tabla alfabéticamente por una columna.
 * Fuente: https://youtu.be/8SL_hM1a0yo
 * TODO: Ordenar campos numéricos correctamente.
 * TODO: Poner campos vacíos al final.
 * 
 * @param {string} tblID El id de la tabla que se va a ordenar
 * @param {number} colIdx El index de la columna por la que ordenar
 */
function sortTableByIdAndColumn(tblID, colIdx) {
	const table = document.getElementById(tblID)
	if (table === null) {
		console.warn("sortable: la tabla ", tblID, " no existe")
		return
	}
	const tbody = table.tBodies[0]
	const rows = Array.from(tbody.querySelectorAll("tr"))

	let ordenASC = true
	if (table.querySelector(`th:nth-child(${ colIdx + 1 })`).classList.contains("th-sort-asc")) {
		ordenASC = false
	}
	
	const sortedRows = rows.sort((a, b) => {
		aColText = normalizar(a.querySelector(`td:nth-child(${ colIdx + 1 })`).textContent)
		bColText = normalizar(b.querySelector(`td:nth-child(${ colIdx + 1 })`).textContent)
		// console.log(bColText + " - " + aColText)
		
		// Opción A: usando matemáticas solamente.
		// const direccion = ordenASC ? 1 : -1
		// return aColText > bColText ? (1 * direccion) : (-1 * direccion)

		// Opción B: usando localCompare para ordenar números correctamente.
		if (ordenASC) {
			return aColText.localeCompare(bColText, undefined, { numeric: true, sensitivity: 'base' });
		} else {
			return bColText.localeCompare(aColText, undefined, { numeric: true, sensitivity: 'base' });
		}
	})

	while (tbody.firstChild) {
		tbody.removeChild(tbody.firstChild)
	}
	tbody.append(...sortedRows)

	// Agregar clase en la celda del encabezado
	table.querySelectorAll("th").forEach(th => th.classList.remove("th-sort-asc", "th-sort-desc"))
	table.querySelector(`th:nth-child(${ colIdx + 1 })`).classList.toggle("th-sort-asc", ordenASC)
	table.querySelector(`th:nth-child(${ colIdx + 1 })`).classList.toggle("th-sort-desc", !ordenASC)
}

// Que todas las tablas se puedan ordenar por las columnas que tengan el atributo "sortable" en sus <th>.
document.querySelectorAll('table').forEach(tbl => {
	if (tbl.id == "") {
		return
	}	
	for(let col of tbl.querySelectorAll("th[sortable]").entries()) {
		col[1].addEventListener("click", () => sortTableByIdAndColumn(tbl.id, col[0]) )
	}
});