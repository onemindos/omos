<?php
	if (isset($_SESSION['UserName']))
		{
?>
			<br>
			<br>
			<table class="tblNavigation">
				<tr>
					<td class="tdNavMenu" valign="middle">
						<input <?=$strDefault;?> class="iButtonNav" onClick="this.disabled='true'; location.href='Default.php'" type="button" value=" Load Defaults ">
					</td>
				</tr>
				<tr>
					<td class="tdNavMenu" valign="middle">
						<input <?=$strSite;?> class="iButtonNav" onClick="this.disabled='true'; location.href='Site.php'" type="button"  value=" Site Configuration ">
					</td>
				</tr>
				<tr>
					<td class="tdNavMenu" valign="middle">
						<input <?=$strGoogle;?> class="iButtonNav" onClick="this.disabled='true'; location.href='Google.php'" type="button"  value=" Adsense Settings ">
					</td>
				</tr>
				<tr>
					<td class="tdNavMenu" valign="middle">
						<input <?=$strAmazon;?> class="iButtonNav" onClick="this.disabled='true'; location.href='Amazon.php'" type="button"  value=" Amazon Settings ">
					</td>
				</tr>
				<tr>
					<td class="tdNavMenu" valign="middle">
						<input <?=$strContact;?> class="iButtonNav" onClick="this.disabled='true'; location.href='Contact.php'" type="button"  value=" Contact Information ">
					</td>
				</tr>
				<tr>
					<td class="tdNavMenu" valign="middle">
						<input <?=$strSecurity;?> class="iButtonNav" onClick="this.disabled='true'; location.href='Security.php'" type="button"  value=" Admin Security ">
					</td>
				</tr>
				<tr>
					<td class="tdNavMenu" valign="middle">
						<input <?=$strPartners;?> class="iButtonNav" onClick="this.disabled='true'; location.href='Partners.php'" type="button"  value=" Partners ">
					</td>
				</tr>
				<tr>
					<td class="tdNavMenu" valign="middle">
						<input <?=$strSitemap;?> class="iButtonNav" onClick="this.disabled='true'; location.href='Sitemap.php'" type="button"  value=" XML Sitemap ">
					</td>
				</tr>				
				<tr>
					<td class="tdNavMenu" valign="middle">
						<input <?=$strSupport;?> class="iButtonNav" onClick="this.disabled='true'; location.href='Support.php'" type="button"  value=" Help Desk ">
					</td>
				</tr>				
				<tr>
					<td class="tdNavMenu" valign="middle">
						<hr size="0">
					</td>
				</tr>
				<tr>
					<td class="tdNavMenu" valign="middle">
						<input class="iButtonNav" onClick="this.disabled='true'; location.href='Session.php?Status=0'" type="button"  value=" Log Out ">
					</td>
				</tr>				
			</table>
<?php
		}
?>		