# Maintainer: Kees van Voorthuizen <me@keesvv.nl>
pkgname=svm
pkgver=0.0.1_pre
_ghver=0.0.1-pre
pkgrel=1
pkgdesc="Lightweight service manager for runit."
arch=('x86_64')
url="https://github.com/keesvv/svm"
license=('GPL')
groups=()
depends=('runit')
makedepends=('go>=1.16')
optdepends=()
provides=()
conflicts=()
replaces=()
backup=()
options=()
install=
changelog=
source=("$pkgname-$pkgver.tar.gz::https://github.com/keesvv/svm/archive/refs/tags/$_ghver.tar.gz")
noextract=()
sha256sums=('a88888e53fd96cf7412784e0e3879c488fcaeab5937ef96ba754ccf553e6be02')

build() {
  cd "$pkgname-$_ghver"
  make
}

package() {
  cd "$pkgname-$_ghver"
  make install
}
