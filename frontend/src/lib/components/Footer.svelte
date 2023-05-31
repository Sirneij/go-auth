<script>
	import Developer from '$lib/img/hero-image.png';

	const year = new Date().getFullYear();
</script>

<footer class="footer-container">
	<div class="footer-branding-container">
		<div class="footer-branding">
			<a class="footer-crafted-by-container" href="https://github.com/Sirneij">
				<span>Developed by</span>
				<img class="footer-branded-crafted-img" src={Developer} alt="John Owolabi Idogun" />
			</a>

			<span class="footer-copyright">&copy; {year} John Owolabi Idogun. All Rights Reserved.</span>
		</div>
	</div>
</footer>
