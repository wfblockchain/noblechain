cd proto
buf generate
cd ..

cp -r github.com/wfblockchain/noblechain/v5/* ./
rm -rf github.com
