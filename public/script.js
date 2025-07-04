document.querySelectorAll('.topbar a').forEach(link => {
  link.addEventListener('click', e => {
    e.preventDefault();
    const page = e.target.getAttribute('data-page');

    // Update topbar
    document.querySelectorAll('.topbar a').forEach(a => a.classList.remove('active'));
    e.target.classList.add('active');

    // Toggle page
    document.querySelectorAll('.page').forEach(section => section.classList.remove('active'));
    document.getElementById(page)?.classList.add('active');
  });
});
