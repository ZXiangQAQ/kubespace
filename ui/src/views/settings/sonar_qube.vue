<template>
    <div>
        <clusterbar :titleName="titleName" :nameFunc="nameSearch" :createFunc="openCreateDialog" createDisplay="添加SonarQube令牌"/>
        <div class="dashboard-container" ref="tableCot">
        <el-table
            ref="multipleTable"
            :data="originSonarQubes"
            class="table-fix"
            :cell-style="cellStyle"
            v-loading="loading"
            :default-sort = "{prop: 'name'}"
            tooltip-effect="dark"
            style="width: 100%"
        >
        <el-table-column prop="name" label="名称" show-overflow-tooltip min-width="15"></el-table-column>
        <el-table-column prop="host_url" label="地址" show-overflow-tooltip min-width="15"></el-table-column>
        <el-table-column prop="update_user" label="操作人" show-overflow-tooltip min-width="10"></el-table-column>
        <el-table-column prop="update_time" label="更新时间" show-overflow-tooltip min-width="15">
          <template slot-scope="scope">
            {{ $dateFormat(scope.row.update_time) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template slot-scope="scope">
            <div class="tableOperate">
              <el-link :disabled="!$editorRole()" :underline="false" type="primary" style="margin-right: 15px;" @click="openUpdateFormDialog(scope.row)">编辑</el-link>
              <el-link :disabled="!$editorRole()" :underline="false" type="danger" @click="handleDeleteSonarQube(scope.row.id, scope.row.name)">删除</el-link>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <el-dialog :title="updateFormVisible ? '修改SonarQube令牌' : '添加SonarQube令牌'" :visible.sync="createFormVisible" @close="closeFormDialog" :destroy-on-close="true">
        <div v-loading="dialogLoading">
          <div class="dialogContent" style="">
            <el-form :model="form" :rules="rules" ref="form" label-position="left" label-width="105px">
              <el-form-item label="名称" prop="name" autofocus>
                <el-input v-model="form.name" autocomplete="off" placeholder="请输入名称，该名称全局唯一" size="small"></el-input>
              </el-form-item>
              <el-form-item label="描述" prop="description">
                <el-input v-model="form.description" type="textarea" autocomplete="off" placeholder="请输入SonarQube令牌描述" size="small"></el-input>
              </el-form-item>
              <el-form-item label="连接地址" prop="host_url">
                <el-input v-model="form.host_url" autocomplete="off" placeholder="请输入SonarQube地址" size="small"></el-input>
              </el-form-item>
              <el-form-item label="令牌" prop="token" :required="true">
                <el-input v-model="form.token" type="password" autocomplete="off" placeholder="请输入令牌" size="small"></el-input>
              </el-form-item>
              <el-alert type="info" title="关于令牌">SonarQube将令牌分成全局令牌、项目令牌、用户令牌，基于权限最小化原则，我们推荐为每个项目配置一个项目级令牌</el-alert>
            </el-form>
          </div>
          <div slot="footer" class="dialogFooter" style="margin-top: 20px;">
            <el-button @click="createFormVisible = false" style="margin-right: 20px;" >取 消</el-button>
            <el-button type="primary" @click="updateFormVisible ? handleUpdateSonarQube() : handleCreateSonarQube()" >确 定</el-button>
          </div>
        </div>
      </el-dialog>
    </div>
    </div>
</template>

<script>
import { Clusterbar } from "@/views/components";
import { createSonarQube, listSonarQube, updateSonarQube, deleteSonarQube } from "@/api/settings/sonar_qube";
import { Message } from "element-ui";

export default {
  name: "SonarQube",
  components: {
    Clusterbar,
  },
  mounted: function () {
    const that = this;
    window.onresize = () => {
      return (() => {
        let heightStyle = window.innerHeight - this.$contentHeight;
        that.maxHeight = heightStyle;
      })();
    };
  },
  created() {
    this.fetchSonarQubes();
  },
  data() {
    return {
      maxHeight: window.innerHeight - this.$contentHeight,
      cellStyle: { border: 0 },
      titleName: ["SonarQube管理"],
      loading: true,
      dialogLoading: false,
      createFormVisible: false,
      updateFormVisible: false,
      form: {
        id: "",
        name: "",
        description: "",
        host_url: "",
        token: "",
      },
      rules: {
        name: [{ required: true, message: '请输入名称', trigger: 'blur' },],
        host_url: [{ required: true, message: '请输入SonarQube地址', trigger: 'blur' },],
        token: [{ required: true, message: '请输入SonarQube令牌', trigger: 'blur' },],
      },
      originSonarQubes: [],
      search_name: "",
    }
  },
  methods: {
    fetchSonarQubes() {
      this.loading = true
      listSonarQube().then((resp) => {
        this.originSonarQubes = resp.data ? resp.data : []
        this.loading = false
      }).catch((err) => {
        console.log(err)
        this.loading = false
      })
    },
    nameSearch(val) {
      this.search_name = val;
    },
    openCreateDialog() {
      this.createFormVisible = true;
    },
    openUpdateFormDialog(row) {
      this.form = {
        id: row.id,
        name: row.name,
        description: row.description,
        host_url: row.host_url,
        token: ''
      }
      this.updateFormVisible = true;
      this.createFormVisible = true;
    },
    closeFormDialog() {
      this.form = {
        id: "",
        name: "",
        description: "",
        host_url: "",
        token: "",
      }
      this.updateFormVisible = false; 
      this.createFormVisible = false;
    },
    handleDeleteSonarQube(id, name) {
      if(!id) {
        Message.error("获取id参数异常，请刷新重试");
        return
      }
      this.$confirm(`请确认是否删除「${name}」此SonarQube配置?`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.loading = true
        deleteSonarQube(id).then(() => {
          this.loading = false
          Message.success("删除SonarQube令牌成功")
          this.fetchSonarQubes()
        }).catch((err) => {
          console.log(err)
          this.loading = false
        });
      }).catch(() => {       
      });
    },
    handleCreateSonarQube() {
      if(!this.form.name) {
        Message.error("SonarQube名称不能为空");
        return
      }
      if(!this.form.host_url) {
        Message.error("请输入SonarQube地址");
        return
      }
      if(!this.form.token) {
        Message.error("请输入SonarQube的令牌");
        return
      }
      let request = {
        name: this.form.name, 
        description: this.form.description, 
        host_url: this.form.host_url,
        token: this.form.token,
      }
      this.dialogLoading = true
      createSonarQube(request).then(() => {
        this.dialogLoading = false
        this.createFormVisible = false;
        Message.success("创建SonarQube令牌成功")
        this.fetchSonarQubes()
      }).catch((err) => {
        this.dialogLoading = false
        console.log(err)
      });
    },
    handleUpdateSonarQube() {
      if(!this.form.id) {
        Message.error("获取SonarQube的id参数异常，请刷新重试");
        return
      }
      let request = {
        name: this.form.name, 
        description: this.form.description, 
        host_url: this.form.host_url,
        token: this.form.token,
      }
      this.dialogLoading = true
      updateSonarQube(this.form.id, request).then(() => {
        this.dialogLoading = false
        this.createFormVisible = false;
        Message.success("更新SonarQube令牌成功")
        this.fetchSonarQubes()
      }).catch((err) => {
        this.dialogLoading = false
        console.log(err)
      });
    },
  }
}
</script>

<style lang="scss" scoped>
@import "~@/styles/variables.scss";
.add-spacelt-p {
  font-family: ui-monospace,SFMono-Regular,SF Mono,Menlo,Consolas,Liberation Mono,monospace
}

</style>