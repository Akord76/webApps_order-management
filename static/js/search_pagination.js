/**
 * Generic Client-side Search and Pagination for Bootstrap Tables
 * @param {Object} config Configuration options
 */
function initTableSearchAndPagination(config) {
  const {
    searchInputId = 'searchInput',
    rowsPerPageId = 'rowsPerPage',
    tableBodyId = 'ordersTableBody',
    rowClass = '.order-row',
    paginationNavId = 'paginationNav',
    pageInfoId = 'pageInfo'
  } = config;

  const searchInput = document.getElementById(searchInputId);
  const rowsPerPageSelect = document.getElementById(rowsPerPageId);
  const tableBody = document.getElementById(tableBodyId);

  if (!tableBody) return; // Exit jika tabel tidak ditemukan di halaman ini

  const allRows = Array.from(tableBody.querySelectorAll(rowClass));
  const paginationNav = document.getElementById(paginationNavId);
  const pageInfo = document.getElementById(pageInfoId);

  let currentPage = 1;
  let rowsPerPage = rowsPerPageSelect ? parseInt(rowsPerPageSelect.value) : 10;
  let filteredRows = [...allRows];

  function filterRows() {
    const query = searchInput ? searchInput.value.toLowerCase().trim() : '';
    filteredRows = allRows.filter(row => row.textContent.toLowerCase().includes(query));
    currentPage = 1;
    renderTable();
  }

  function renderTable() {
    allRows.forEach(row => (row.style.display = 'none'));

    const totalRows = filteredRows.length;
    const totalPages = Math.ceil(totalRows / rowsPerPage) || 1;

    if (currentPage > totalPages) currentPage = totalPages;

    const start = (currentPage - 1) * rowsPerPage;
    const end = start + rowsPerPage;
    const visibleRows = filteredRows.slice(start, end);

    visibleRows.forEach(row => (row.style.display = ''));

    if (pageInfo) {
      if (totalRows === 0) {
        pageInfo.textContent = 'Showing 0 entries';
      } else {
        pageInfo.textContent = `Showing ${start + 1} to ${Math.min(end, totalRows)} of ${totalRows} entries`;
      }
    }

    renderPagination(totalPages);
  }

  function renderPagination(totalPages) {
    if (!paginationNav) return;
    paginationNav.innerHTML = '';

    if (totalPages <= 1) return;

    // Previous Button
    const prevLi = document.createElement('li');
    prevLi.className = `page-item ${currentPage === 1 ? 'disabled' : ''}`;
    prevLi.innerHTML = `<a class="page-link" href="#">Previous</a>`;
    prevLi.addEventListener('click', (e) => {
      e.preventDefault();
      if (currentPage > 1) {
        currentPage--;
        renderTable();
      }
    });
    paginationNav.appendChild(prevLi);

    // Page Number Buttons
    for (let i = 1; i <= totalPages; i++) {
      const pageLi = document.createElement('li');
      pageLi.className = `page-item ${i === currentPage ? 'active' : ''}`;
      pageLi.innerHTML = `<a class="page-link" href="#">${i}</a>`;
      pageLi.addEventListener('click', (e) => {
        e.preventDefault();
        currentPage = i;
        renderTable();
      });
      paginationNav.appendChild(pageLi);
    }

    // Next Button
    const nextLi = document.createElement('li');
    nextLi.className = `page-item ${currentPage === totalPages ? 'disabled' : ''}`;
    nextLi.innerHTML = `<a class="page-link" href="#">Next</a>`;
    nextLi.addEventListener('click', (e) => {
      e.preventDefault();
      if (currentPage < totalPages) {
        currentPage++;
        renderTable();
      }
    });
    paginationNav.appendChild(nextLi);
  }

  // Event Listeners
  if (searchInput) searchInput.addEventListener('input', filterRows);
  if (rowsPerPageSelect) {
    rowsPerPageSelect.addEventListener('change', function () {
      rowsPerPage = parseInt(this.value);
      currentPage = 1;
      renderTable();
    });
  }

  // Initial Render
  renderTable();
}