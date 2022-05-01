#!/bin/bash

npx babel js_src --out-dir static/js --presets react-app/prod --minified --watch --ignore=m_sdk.js
