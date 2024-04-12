<template>
  <div v-loading="loading">
    <el-form :model="params" ref="job" label-position="left" label-width="120px">
        <el-form-item label="扫描镜像" prop="" :required="true">
          <el-select v-model="params.sonar_scanner_image" placeholder="请选择扫描镜像" size="small" style="width: 400px">
            <el-option
              v-for="res in imageResources"
              :key="res.id"
              :label="res.name"
              :value="res.id">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="扫描方式" prop="sonar_scanner_type">
            <el-radio-group v-model="params.sonar_scanner_type">
                <el-radio label="file">配置文件</el-radio>
                <el-radio label="script">自定义配置</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="配置文件" prop="" v-if="params.sonar_scanner_type == 'file'">
          <el-input style="width: 400px;" v-model="params.sonar_scanner_file" autocomplete="off" size="small"
            placeholder="默认为当前代码库根目录下的「sonar-project.properties」文件">
          </el-input>
        </el-form-item>
        <el-form-item label="自定义配置" prop="" v-if="params.sonar_scanner_type == 'script'">
          <el-input type="textarea" :rows="6" v-model="params.sonar_scanner_script" autocomplete="off" size="small" 
            placeholder="# must be unique in a given SonarQube instance
sonar.projectKey="></el-input>  
        </el-form-item>
        <el-form-item label="SonarQube令牌" prop="">
          <el-select v-model="params.sonar_qube_id" placeholder="请选择SonarQube令牌" size="small" style="width: 400px;">
            <el-option 
              v-for="item in sonarQubes"
              :key="item.id"
              :label="item.name"
              :value="item.id">
            </el-option>
          </el-select>
        </el-form-item>
    </el-form>
  </div>
</template>

<script>
import { listResources } from "@/api/pipeline/resource"
import { listSonarQube } from "@/api/settings/sonar_qube"

export default {
    name: 'SonarQube',
    data() {
      return {
        loading: false,
        resources: [],
        resource: this.params.resource ? this.params.resource : 0,
        sonarQubes: [],
      }
    },
    props: ['params'],
    computed: {
      workspaceId() {
        return this.$route.params.workspaceId
      },
      imageResources() {
        let res = []
        for(let r of this.resources) {
          if(r.type == 'image') {
            res.push(r)
          }
        }
        return res
      },
    },
    beforeMount() {
        if(!this.params.sonar_scanner_type) {
          this.$set(this.params, 'sonar_scanner_type', 'file')
        }
        this.fetchResources()
        this.fetchSonarQubes()
    },
    methods: {
      fetchResources() {
        this.loading = true
        listResources(this.workspaceId).then((resp) => {
            this.resources = resp.data ? resp.data : []
            this.loading = false
        }).catch((err) => {
            console.log(err)
            this.loading = false
        })
      },
      fetchSonarQubes() {
        this.loading = true
        listSonarQube().then((resp) => {
            this.sonarQubes = resp.data ? resp.data : []
            this.loading = false
        }).catch((err) => {
            console.log(err)
            this.loading = false
        })
      }
    }
}
</script>

