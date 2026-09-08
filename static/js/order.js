// static/js/order.js

// -------------------------------------------------------------
// Helper Global untuk API Call
// -------------------------------------------------------------
async function apiCall(url, method = 'GET', body = null) {
    const options = {
        method: method,
        headers: {
            'Content-Type': 'application/json'
        }
    };

    if (body) {
        options.body = JSON.stringify(body);
    }

    try {
        const response = await fetch(url, options);
        if (!response.ok) {
            const errData = await response.json();
            throw new Error(errData.error || 'Terjadi kesalahan pada request');
        }
        return await response.json();
    } catch (error) {
        console.error('API Error:', error);
        alert(error.message);
    }
}

// -------------------------------------------------------------
// Fungsi-fungsi Order Master
// -------------------------------------------------------------
async function loadAllOrders() {
    const orders = await apiCall('/orders');
    console.log("Daftar Order:", orders);
    return orders;
}

async function getOrderDetail(orderID, orderNo) {
    const url = `/orders/${orderID}?orderNo=${encodeURIComponent(orderNo)}`;
    const order = await apiCall(url, 'GET');
    console.log("Detail Order Master:", order);
    return order;
}

async function createOrder(payload) {
    const response = await apiCall('/orders', 'POST', payload);
    if (response) {
        alert("Order berhasil dibuat!");
        window.location.reload();
    }
}

async function updateOrder(orderID, orderNo, payload) {
    const url = `/orders/${orderID}?orderNo=${encodeURIComponent(orderNo)}`;
    const response = await apiCall(url, 'PUT', payload);
    if (response) {
        alert("Order berhasil diperbarui!");
        window.location.reload();
    }
}

async function deleteOrder(orderID, orderNo) {
    if (!confirm("Apakah Anda yakin ingin menghapus order ini?")) return;

    const url = `/orders/${orderID}?orderNo=${encodeURIComponent(orderNo)}`;
    const response = await fetch(url, { method: 'DELETE' });
    
    if (response.ok) {
        alert("Order berhasil dihapus!");
        window.location.reload();
    } else {
        alert("Gagal menghapus order");
    }
}

// -------------------------------------------------------------
// Fungsi-fungsi Order Detail
// -------------------------------------------------------------
async function addOrderDetail(orderNo, payload) {
    const url = `/orders/details?orderNo=${encodeURIComponent(orderNo)}`;
    const response = await apiCall(url, 'POST', payload);
    if (response) {
        alert("Item detail berhasil ditambahkan!");
        window.location.reload();
    }
}

async function fetchOrderDetails(orderID, orderNo) {
    const url = `/orders/${orderID}/details?orderNo=${encodeURIComponent(orderNo)}`;
    const details = await apiCall(url, 'GET');
    return details;
}

async function deleteOrderDetail(detailID, orderNo) {
    if (!confirm("Hapus item detail ini?")) return;

    const url = `/orders/details/${detailID}?orderNo=${encodeURIComponent(orderNo)}`;
    const response = await fetch(url, { method: 'DELETE' });
    
    if (response.ok) {
        alert("Detail item berhasil dihapus!");
        window.location.reload();
    }
}

// -------------------------------------------------------------
// Fungsi-fungsi Autocomplete
// -------------------------------------------------------------
async function getAutocompleteData() {
    return await apiCall('/orders/autocomplete');
}

async function getOrderAutocompleteByCustomer(custID, orderNoInput) {
    const url = `/orders/autocomplete/${encodeURIComponent(custID)}?orderNo=${encodeURIComponent(orderNoInput)}`;
    return await apiCall(url, 'GET');
}