
## Extract parameters from ConfigMaps and Secrets as files

mkdir params
cd params
mkdir ConfigMap__v1_holdings_holdings-api-config.yaml
oc extract -f ../output/resources/holdings/ConfigMap__v1_holdings_holdings-api-config.yaml --to ConfigMap__v1_holdings_holdings-api-config.yaml

```bash
cd myproject/source/holdings-versionchanged-parameterized/kustomize/base
mv holdings-api-config-configmap.yaml holdings-api-config-configmap.yaml.bak
sed -i.bak '/holdings-api-config-configmap.yaml/d' kustomization.yaml
...

3```bash
echo -n "configMapGenerator:" >> kustomization.yaml
echo -n "
- name: holdings-api-config
  behavior: create
  files:" >> kustomization.yaml

for f in params/ConfigMap__v1_holdings_holdings-api-config.yaml/*
do
  echo -n "
  - $f" >> kustomization.yaml
done

echo -n "
  options:
    disableNameSuffixHash: true
" >> kustomization.yaml
```

In overlays use envs instead:
```yaml
configMapGenerator:
  - name: holdings-api-config
    behavior: merge
    envs:
      - custom.env
```

**TODO**: different file naming conventions



## Remove sensitive data from defaults
## Remove ConfigMaps and Secrets
* If not removed from exported files, they would be created twice
* On the other hand, if we only use the configmapGenerator, we'd lose the labels and annotations from the original reasource

## Use configMapGenerator and secretGenerator to create ConfigMaps and Secrets from files
## Define hooks for overriding the defaults (for all sensitive data as well)
## Encryption
## Howto Helm?
## Integrate in move2kube

## Conclusion
* Solution doesn't scale
  * Not flexible, lots of limitations
  * Ad-hoc implementation is preferrable
    * To keep/filter/manipulate original labels/annotations
    * To customize installation options
  * With Konveyor tools, only per-namespace implementation is doable 
    * Complex manipulation of target namespaces




