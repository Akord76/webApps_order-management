(function () {
    'use strict';

    /**
 * Autocomplete untuk Products
 */
function initProductAutocomplete(options = {}) {
    const config = {
        searchUrl: options.searchUrl || '/commitmentfees/products/search',
        minChars: options.minChars || 2,
        debounceMs: options.debounceMs || 300,
        inputContainerId: options.inputId || 'productName',
        productIdField: options.productIdId || 'productId',
        measureField: options.measureId || 'measure'
    };

    const input = document.getElementById(config.inputContainerId);
    if (!input) return;

    let debounceTimer;
    let selectedIndex = -1;

    const list = createDropdown(input);

    input.addEventListener('input', (e) => {
        clearTimeout(debounceTimer);
        const q = e.target.value.trim();
        if (q.length < config.minChars) {
            hideSuggestions(list);
            return;
        }
        debounceTimer = setTimeout(() => fetchProductSuggestions(q, list, config), config.debounceMs);
    });

    input.addEventListener('keydown', (e) => handleKeydown(e, list));

    document.addEventListener('click', (e) => {
        if (!input.contains(e.target) && !list.contains(e.target)) {
            hideSuggestions(list);
        }
    });

    function handleKeydown(e, dropdownList) {
        const items = dropdownList.querySelectorAll('.list-group-item');
        if (!items.length) return;

        if (e.key === 'ArrowDown') {
            e.preventDefault();
            selectedIndex = Math.min(selectedIndex + 1, items.length - 1);
            updateActive(items);
        } else if (e.key === 'ArrowUp') {
            e.preventDefault();
            selectedIndex = Math.max(selectedIndex - 1, 0);
            updateActive(items);
        } else if (e.key === 'Enter' && selectedIndex >= 0) {
            e.preventDefault();
            items[selectedIndex].click();
        } else if (e.key === 'Escape') {
            hideSuggestions(dropdownList);
        }
    }

    function updateActive(items) {
        items.forEach(i => i.classList.remove('active'));
        if (selectedIndex >= 0) {
            items[selectedIndex].classList.add('active');
            items[selectedIndex].scrollIntoView({ block: 'nearest' });
        }
    }

    function hideSuggestions(dropdownList) {
        dropdownList.style.display = 'none';
        dropdownList.innerHTML = '';
        selectedIndex = -1;
    }

async function fetchProductSuggestions(q, dropdownList, cfg) {
    // try {
        const token = localStorage.getItem('token') || sessionStorage.getItem('token');
        const res = await fetch(`${cfg.searchUrl}?q=${encodeURIComponent(q.trim())}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                ...(token && { 'Authorization': `Bearer ${token}` })
            }
        });

        // Jika HTTP Status bukan 200 OK
        if (!res.ok) {
            const errText = await res.text();
            console.error(`[API Error ${res.status}]:`, errText);
            hideSuggestions(dropdownList);
            return;
        }

        const responseData = await res.json();
        console.log("Response JSON dari server:", responseData);

        // Ambil array dari key 'results' (sesuai Postman) atau 'data'
        const productList = responseData.results || responseData.data || [];

        renderProductSuggestions(productList, dropdownList);
    // } catch (err) {
    //     console.error("Gagal melakukan fetch autocomplete produk:", err);
    //     hideSuggestions(dropdownList);
    // }
}
    function renderProductSuggestions(products, dropdownList) {
        dropdownList.innerHTML = '';
        selectedIndex = -1;

        if (!products || !Array.isArray(products) || products.length === 0) {
            dropdownList.innerHTML = '<li class="list-group-item text-muted">Tidak ditemukan</li>';
            dropdownList.style.display = 'block';
            return;
        }

        products.forEach(p => {
            const li = document.createElement('li');
            li.className = 'list-group-item list-group-item-action';

            // PERBAIKAN 2: Fallback property name (CamelCase & snake_case)
            const name = p.ProductName || p.product_name || p.productName || '';
            const id = p.ProductID || p.product_id || p.productId || '';
            const measure = p.Measure || p.measure || '';

            li.textContent = name;
            li.dataset.id = id;
            li.dataset.measure = measure;

            li.onclick = () => selectProduct({ ProductName: name, ProductID: id, Measure: measure }, dropdownList);
            dropdownList.appendChild(li);
        });
        dropdownList.style.display = 'block';
    }

    function selectProduct(p, dropdownList) {
        document.getElementById(config.inputContainerId).value = p.ProductName;
        const targetId = document.getElementById(config.productIdField);
        const targetMeasure = document.getElementById(config.measureField);
        if (targetId) targetId.value = p.ProductID;
        if (targetMeasure) targetMeasure.value = p.Measure;
        hideSuggestions(dropdownList);
    }

    function createDropdown(targetInput) {
        let list = document.getElementById('suggestionsList');
        if (!list) {
            list = document.createElement('ul');
            list.id = 'suggestionsList';
            list.className = 'list-group position-absolute w-100 shadow-sm';
            list.style.cssText = 'z-index: 1000; display: none; max-height: 250px; overflow-y: auto; border-radius: 0 0 0.375rem 0.375rem;';
            targetInput.parentElement.classList.add('position-relative');
            targetInput.parentElement.appendChild(list);
        }
        return list;
    }
}

    /**
     * Autocomplete untuk Karyawan
     */
    function initEmployeeAutocomplete(options = {}) {
        const config = {
            searchUrl: options.searchUrl || '/commitmentfees/employees/search',
            inputId: options.inputId || 'employee_search',
            listId: options.listId || 'suggestions',
            hiddenCardId: options.hiddenCardId || 'employee_card_number',
            displayCardId: options.displayCardId || 'display_card_number',
            debounceMs: options.debounceMs || 300
        };

        const searchInput = document.getElementById(config.inputId);
        const suggestionsList = document.getElementById(config.listId);
        const hiddenCardInput = document.getElementById(config.hiddenCardId);
        const displayCardInput = document.getElementById(config.displayCardId);

        if (!searchInput || !suggestionsList) return;

        let debounceTimer;

        searchInput.addEventListener('input', function () {
            const query = this.value.trim();

            if (query.length < 2) {
                suggestionsList.style.display = 'none';
                suggestionsList.innerHTML = '';
                if (hiddenCardInput) hiddenCardInput.value = '';
                if (displayCardInput) displayCardInput.value = '';
                return;
            }

            clearTimeout(debounceTimer);
            debounceTimer = setTimeout(() => {
                fetchEmployees(query);
            }, config.debounceMs);
        });

        function fetchEmployees(query) {
            fetch(`${config.searchUrl}?q=${encodeURIComponent(query)}`)
                .then(response => response.json())
                .then(res => {
                    suggestionsList.innerHTML = '';

                    if (res.success && res.data && res.data.length > 0) {
                        res.data.forEach(item => {
                            const li = document.createElement('li');
                            li.className = 'list-group-item list-group-item-action';
                            const fullName = `${item.first_name} ${item.last_name}`;

                            li.innerHTML = `
                                <div><strong>${fullName}</strong></div>
                                <small class="text-muted">${item.email}</small>
                            `;

                            li.addEventListener('click', function () {
                                searchInput.value = fullName;
                                if (hiddenCardInput) hiddenCardInput.value = item.employee_card_number;
                                if (displayCardInput) displayCardInput.value = item.employee_card_number;
                                suggestionsList.style.display = 'none';
                            });

                            suggestionsList.appendChild(li);
                        });
                        suggestionsList.style.display = 'block';
                    } else {
                        suggestionsList.style.display = 'none';
                    }
                })
                .catch(err => {
                    console.error('Error fetching autocomplete:', err);
                    suggestionsList.style.display = 'none';
                });
        }

        document.addEventListener('click', function (e) {
            if (!searchInput.contains(e.target) && !suggestionsList.contains(e.target)) {
                suggestionsList.style.display = 'none';
            }
        });
    }

    // Ekspor fungsi ke scope window agar bisa dipanggil dari HTML Go
    window.AppAutocomplete = {
        initProduct: initProductAutocomplete,
        initEmployee: initEmployeeAutocomplete
    };
})();