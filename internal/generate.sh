#!/bin/bash
#
# Wraith Game Engine
# Copyright (c) 2022 Michael D. Henderson
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published
# by the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.
#


for dfn in definition/*.go; do
  pkg="${dfn#definition/}"; pkg="${pkg%.go}"
  echo " info: generating package for ${pkg}..."
  mkdir -p services/${pkg}
  out="services/${pkg}/oto-service.go"
  ~/go/bin/oto \
    -template otohttp/templates/oto-server.go.plush \
    -out "${out}" \
    -pkg "${pkg}" \
    "${dfn}"
  gofmt -w "${out}" "${out}"
done

exit 0

#~/go/bin/oto -template otohttp/templates/oto-client.go.plush \
#  -out client/oto-client.go \
#  -ignore Ignorer \
#  -pkg main \
#  definition
#
#gofmt -w client/oto-client.go client/oto-client.go
